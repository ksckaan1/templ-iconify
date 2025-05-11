package service

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"slices"
	"strings"
	"unicode"

	"github.com/ksckaan1/templ-iconify/internal/core/customerrors"
	"github.com/ksckaan1/templ-iconify/internal/core/domain"
	"github.com/ksckaan1/templ-iconify/internal/infra/iconifyclient"
	"github.com/samber/lo"
)

type IconService struct {
	iconifyClient *iconifyclient.IconifyClient
}

func NewIconService(iconifyClient *iconifyclient.IconifyClient) *IconService {
	return &IconService{
		iconifyClient: iconifyClient,
	}
}

func (s *IconService) FindIcons(ctx context.Context, names ...string) ([]*domain.Icon, error) {
	var targetIcons []*domain.Icon

	for _, name := range names {
		icon, err := s.parseIconName(name)
		if err != nil {
			return nil, customerrors.ErrInvalidIconName
		}

		icons, err := s.generateDownloadList(ctx, icon)
		if err != nil {
			return nil, fmt.Errorf("generateDownloadList: %w", err)
		}
		targetIcons = append(targetIcons, icons...)
	}

	targetIcons = lo.Uniq(targetIcons)

	if len(targetIcons) == 0 {
		return nil, customerrors.ErrIconNotFound
	}

	slices.SortFunc(targetIcons, func(a, b *domain.Icon) int {
		return strings.Compare(a.Prefix+":"+a.Name, b.Prefix+":"+b.Name)
	})

	return targetIcons, nil
}

func (s *IconService) DownloadIcon(ctx context.Context, icon *domain.Icon, saveDir string) error {
	svgBody, err := s.iconifyClient.GetIconSVG(ctx, icon.Prefix, icon.Name)
	if err != nil {
		return fmt.Errorf("iconifyClient.GetIconSVG: %w", err)
	}
	icon.SVGBody = svgBody
	err = s.generateTemplate(icon)
	if err != nil {
		return fmt.Errorf("generateTemplate: %w", err)
	}
	err = s.saveTemplate(icon, saveDir)
	if err != nil {
		return fmt.Errorf("saveTemplate: %w", err)
	}
	return nil
}

func (s *IconService) generateDownloadList(ctx context.Context, icon *domain.Icon) ([]*domain.Icon, error) {
	collections, err := s.iconifyClient.GetCollections(ctx)
	if err != nil {
		return nil, fmt.Errorf("iconifyClient.GetCollections: %w", err)
	}

	filteredCollections, err := s.filterCollections(icon.Prefix, collections)
	if err != nil {
		return nil, fmt.Errorf("filterCollections: %w", err)
	}

	iconWildcard := s.wildcardToRegex(icon.Name)
	rgxIcon, err := regexp.Compile(iconWildcard)
	if err != nil {
		return nil, fmt.Errorf("regexp.Compile: %w", err)
	}

	var result []*domain.Icon

	for _, collection := range filteredCollections {
		icons, err := s.iconifyClient.GetCollectionIcons(ctx, collection)
		if err != nil {
			return nil, fmt.Errorf("iconifyClient.GetCollectionIcons: %w", err)
		}

		filteredIcons, err := s.filterIcons(rgxIcon, icons)
		if err != nil {
			return nil, fmt.Errorf("filterIcons: %w", err)
		}

		result = append(result, lo.Map(filteredIcons, func(icon string, _ int) *domain.Icon {
			return &domain.Icon{
				Prefix: collection,
				Name:   icon,
			}
		})...)
	}

	return result, nil
}

func (s *IconService) filterCollections(targetCollection string, allCollections []string) ([]string, error) {
	rgxCollection, err := regexp.Compile(s.wildcardToRegex(targetCollection))
	if err != nil {
		return nil, fmt.Errorf("regexp.Compile: %w", err)
	}
	return lo.Filter(allCollections, func(collection string, _ int) bool {
		return rgxCollection.MatchString(collection)
	}), nil
}

func (s *IconService) filterIcons(rgxIcon *regexp.Regexp, allIcons []string) ([]string, error) {
	return lo.Filter(allIcons, func(icon string, _ int) bool {
		return rgxIcon.MatchString(icon)
	}), nil
}

var rgxIconName = regexp.MustCompile(`^([\w\-\*]+):([\w\-\*]+)$`)

func (s *IconService) parseIconName(name string) (*domain.Icon, error) {
	matches := rgxIconName.FindStringSubmatch(name)
	if len(matches) != 3 {
		return nil, errors.New("invalid icon name")
	}
	return &domain.Icon{
		Prefix: matches[1],
		Name:   matches[2],
	}, nil
}

func (s *IconService) wildcardToRegex(pattern string) string {
	anchoredPrefix := ""
	if !strings.HasPrefix(pattern, "*") {
		anchoredPrefix = "^"
	}
	anchoredSuffix := ""
	if !strings.HasSuffix(pattern, "*") {
		anchoredSuffix = "$"
	}
	parts := strings.Split(pattern, "*")
	var regexBuilder strings.Builder
	regexBuilder.WriteString(anchoredPrefix)
	for i, part := range parts {
		if i > 0 {
			regexBuilder.WriteString(".*")
		}
		regexBuilder.WriteString(regexp.QuoteMeta(part))
	}
	regexBuilder.WriteString(anchoredSuffix)
	return regexBuilder.String()
}

const iconTemplate = `package %s

type %sProps struct {
	Width  string
	Height string
	Color  string
}

templ %s(props %sProps) {
	%s
}
`

func (s *IconService) generateTemplate(icon *domain.Icon) error {
	icon.PkgName = s.generatePackageName(icon.Prefix)
	icon.ComponentName = s.generateComponentName(icon.Name)
	icon.TemplateString = fmt.Sprintf(iconTemplate, icon.PkgName, icon.ComponentName, icon.ComponentName, icon.ComponentName, s.replaceSVG(icon.SVGBody))
	return nil
}

func (s *IconService) generatePackageName(prefix string) string {
	pkgName := strings.ToLower(prefix)
	pkgName = strings.ReplaceAll(pkgName, "-", "")
	if unicode.IsDigit(rune(pkgName[0])) {
		pkgName = "icon" + pkgName
	}
	return pkgName
}

func (s *IconService) generateComponentName(iconName string) string {
	parts := strings.Split(iconName, "-")
	for i, part := range parts {
		if len(part) > 0 {
			parts[i] = strings.ToUpper(part[:1]) + part[1:]
		}
	}
	result := strings.Join(parts, "")
	if unicode.IsDigit(rune(result[0])) {
		result = "Icon" + result
	}
	return result
}

var rgxSizeReplacer = regexp.MustCompile(`width="[\w]+" height="[\w]+"`)

const sizeReplacerTarget = `
		if props.Width != "" {
			width={ props.Width }
		}
		if props.Height != "" {
			height={ props.Height }
		}
		`
const colorReplacerTarget = `
			if props.Color != "" {
				fill={ props.Color }
			} else {
				fill="currentColor"
			}
			`

func (s *IconService) replaceSVG(svgBody string) string {
	out := rgxSizeReplacer.ReplaceAllString(svgBody, sizeReplacerTarget)
	out = strings.ReplaceAll(out, `fill="currentColor"`, colorReplacerTarget)
	return out
}

func (s *IconService) saveTemplate(icon *domain.Icon, saveDir string) error {
	err := os.MkdirAll(filepath.Join(saveDir, icon.PkgName), 0755)
	if err != nil {
		return fmt.Errorf("os.MkdirAll: %w", err)
	}
	savePath := filepath.Join(saveDir, icon.PkgName, icon.Name+".templ")
	file, err := os.Create(savePath)
	if err != nil {
		return fmt.Errorf("os.Create: %w", err)
	}
	defer file.Close()
	_, err = file.WriteString(icon.TemplateString)
	if err != nil {
		return fmt.Errorf("file.WriteString: %w", err)
	}
	return nil
}
