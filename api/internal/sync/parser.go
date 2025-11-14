package sync

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

// TemplateParser parses template manifests from repositories
type TemplateParser struct{}

// NewTemplateParser creates a new template parser
func NewTemplateParser() *TemplateParser {
	return &TemplateParser{}
}

// ParsedTemplate represents a parsed template from a repository
type ParsedTemplate struct {
	Name        string
	DisplayName string
	Description string
	Category    string
	AppType     string
	Icon        string
	Manifest    string   // JSON-encoded full manifest
	Tags        []string
}

// TemplateManifest represents the YAML structure of a template file
type TemplateManifest struct {
	APIVersion string `yaml:"apiVersion"`
	Kind       string `yaml:"kind"`
	Metadata   struct {
		Name      string            `yaml:"name"`
		Namespace string            `yaml:"namespace,omitempty"`
		Labels    map[string]string `yaml:"labels,omitempty"`
	} `yaml:"metadata"`
	Spec struct {
		DisplayName      string            `yaml:"displayName"`
		Description      string            `yaml:"description"`
		Category         string            `yaml:"category"`
		AppType          string            `yaml:"appType,omitempty"`
		Icon             string            `yaml:"icon,omitempty"`
		BaseImage        string            `yaml:"baseImage"`
		DefaultResources map[string]string `yaml:"defaultResources,omitempty"`
		Ports            []struct {
			Name          string `yaml:"name"`
			ContainerPort int    `yaml:"containerPort"`
			Protocol      string `yaml:"protocol,omitempty"`
		} `yaml:"ports,omitempty"`
		Env          []map[string]interface{} `yaml:"env,omitempty"`
		VolumeMounts []map[string]interface{} `yaml:"volumeMounts,omitempty"`
		VNC          *struct {
			Enabled    bool   `yaml:"enabled"`
			Port       int    `yaml:"port"`
			Protocol   string `yaml:"protocol,omitempty"`
			Encryption bool   `yaml:"encryption,omitempty"`
		} `yaml:"vnc,omitempty"`
		WebApp *struct {
			Enabled     bool   `yaml:"enabled"`
			Port        int    `yaml:"port"`
			Path        string `yaml:"path,omitempty"`
			HealthCheck string `yaml:"healthCheck,omitempty"`
		} `yaml:"webapp,omitempty"`
		Capabilities []string `yaml:"capabilities,omitempty"`
		Tags         []string `yaml:"tags,omitempty"`
	} `yaml:"spec"`
}

// ParseRepository parses all template manifests in a repository
func (p *TemplateParser) ParseRepository(repoPath string) ([]*ParsedTemplate, error) {
	var templates []*ParsedTemplate

	// Walk through repository looking for YAML files
	err := filepath.WalkDir(repoPath, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		// Skip directories
		if d.IsDir() {
			// Skip .git directory
			if d.Name() == ".git" {
				return filepath.SkipDir
			}
			return nil
		}

		// Only process YAML files
		ext := strings.ToLower(filepath.Ext(path))
		if ext != ".yaml" && ext != ".yml" {
			return nil
		}

		// Parse template file
		template, err := p.ParseTemplateFile(path)
		if err != nil {
			// Log error but continue processing other files
			// (not all YAML files may be templates)
			return nil
		}

		templates = append(templates, template)
		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to walk repository: %w", err)
	}

	return templates, nil
}

// ParseTemplateFile parses a single template YAML file
func (p *TemplateParser) ParseTemplateFile(filePath string) (*ParsedTemplate, error) {
	// Read file
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	// Parse YAML
	var manifest TemplateManifest
	if err := yaml.Unmarshal(data, &manifest); err != nil {
		return nil, fmt.Errorf("failed to parse YAML: %w", err)
	}

	// Validate this is a Template resource
	if manifest.Kind != "Template" {
		return nil, fmt.Errorf("not a Template resource (kind: %s)", manifest.Kind)
	}

	if manifest.APIVersion != "stream.streamspace.io/v1alpha1" {
		return nil, fmt.Errorf("unsupported API version: %s", manifest.APIVersion)
	}

	// Validate required fields
	if manifest.Metadata.Name == "" {
		return nil, fmt.Errorf("template name is required")
	}

	if manifest.Spec.DisplayName == "" {
		return nil, fmt.Errorf("displayName is required")
	}

	if manifest.Spec.BaseImage == "" {
		return nil, fmt.Errorf("baseImage is required")
	}

	// Determine app type
	appType := manifest.Spec.AppType
	if appType == "" {
		// Infer from VNC/WebApp config
		if manifest.Spec.WebApp != nil && manifest.Spec.WebApp.Enabled {
			appType = "webapp"
		} else {
			appType = "desktop"
		}
	}

	// Convert full manifest to JSON for storage
	manifestJSON, err := json.Marshal(manifest)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal manifest to JSON: %w", err)
	}

	// Create parsed template
	template := &ParsedTemplate{
		Name:        manifest.Metadata.Name,
		DisplayName: manifest.Spec.DisplayName,
		Description: manifest.Spec.Description,
		Category:    manifest.Spec.Category,
		AppType:     appType,
		Icon:        manifest.Spec.Icon,
		Manifest:    string(manifestJSON),
		Tags:        manifest.Spec.Tags,
	}

	// Default empty tags to empty array
	if template.Tags == nil {
		template.Tags = []string{}
	}

	return template, nil
}

// ParseTemplateFromString parses a template from a YAML string
func (p *TemplateParser) ParseTemplateFromString(yamlContent string) (*ParsedTemplate, error) {
	// Parse YAML
	var manifest TemplateManifest
	if err := yaml.Unmarshal([]byte(yamlContent), &manifest); err != nil {
		return nil, fmt.Errorf("failed to parse YAML: %w", err)
	}

	// Validate this is a Template resource
	if manifest.Kind != "Template" {
		return nil, fmt.Errorf("not a Template resource (kind: %s)", manifest.Kind)
	}

	// Determine app type
	appType := manifest.Spec.AppType
	if appType == "" {
		if manifest.Spec.WebApp != nil && manifest.Spec.WebApp.Enabled {
			appType = "webapp"
		} else {
			appType = "desktop"
		}
	}

	// Convert full manifest to JSON for storage
	manifestJSON, err := json.Marshal(manifest)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal manifest to JSON: %w", err)
	}

	template := &ParsedTemplate{
		Name:        manifest.Metadata.Name,
		DisplayName: manifest.Spec.DisplayName,
		Description: manifest.Spec.Description,
		Category:    manifest.Spec.Category,
		AppType:     appType,
		Icon:        manifest.Spec.Icon,
		Manifest:    string(manifestJSON),
		Tags:        manifest.Spec.Tags,
	}

	if template.Tags == nil {
		template.Tags = []string{}
	}

	return template, nil
}

// ValidateTemplateManifest validates a template manifest structure
func (p *TemplateParser) ValidateTemplateManifest(yamlContent string) error {
	var manifest TemplateManifest
	if err := yaml.Unmarshal([]byte(yamlContent), &manifest); err != nil {
		return fmt.Errorf("invalid YAML: %w", err)
	}

	// Check required fields
	if manifest.Kind != "Template" {
		return fmt.Errorf("kind must be 'Template', got '%s'", manifest.Kind)
	}

	if manifest.APIVersion != "stream.streamspace.io/v1alpha1" {
		return fmt.Errorf("apiVersion must be 'stream.streamspace.io/v1alpha1', got '%s'", manifest.APIVersion)
	}

	if manifest.Metadata.Name == "" {
		return fmt.Errorf("metadata.name is required")
	}

	if manifest.Spec.DisplayName == "" {
		return fmt.Errorf("spec.displayName is required")
	}

	if manifest.Spec.BaseImage == "" {
		return fmt.Errorf("spec.baseImage is required")
	}

	// Validate app type if specified
	if manifest.Spec.AppType != "" && manifest.Spec.AppType != "desktop" && manifest.Spec.AppType != "webapp" {
		return fmt.Errorf("spec.appType must be 'desktop' or 'webapp', got '%s'", manifest.Spec.AppType)
	}

	return nil
}

// ========== Plugin Parsing ==========

// PluginParser parses plugin manifests from repositories
type PluginParser struct{}

// NewPluginParser creates a new plugin parser
func NewPluginParser() *PluginParser {
	return &PluginParser{}
}

// ParsedPlugin represents a parsed plugin from a repository
type ParsedPlugin struct {
	Name        string
	Version     string
	DisplayName string
	Description string
	Category    string
	PluginType  string
	Icon        string
	Manifest    string   // JSON-encoded full manifest
	Tags        []string
}

// PluginManifest represents the structure of a plugin manifest.json file
type PluginManifest struct {
	Name          string            `json:"name"`
	Version       string            `json:"version"`
	DisplayName   string            `json:"displayName"`
	Description   string            `json:"description"`
	Author        string            `json:"author"`
	Homepage      string            `json:"homepage,omitempty"`
	Repository    string            `json:"repository,omitempty"`
	License       string            `json:"license,omitempty"`
	Type          string            `json:"type"`
	Category      string            `json:"category,omitempty"`
	Tags          []string          `json:"tags,omitempty"`
	Icon          string            `json:"icon,omitempty"`
	Requirements  map[string]string `json:"requirements,omitempty"`
	Entrypoints   map[string]string `json:"entrypoints,omitempty"`
	ConfigSchema  map[string]interface{} `json:"configSchema,omitempty"`
	DefaultConfig map[string]interface{} `json:"defaultConfig,omitempty"`
	Permissions   []string          `json:"permissions,omitempty"`
	Dependencies  map[string]string `json:"dependencies,omitempty"`
}

// ParseRepository parses all plugin manifests in a repository
func (p *PluginParser) ParseRepository(repoPath string) ([]*ParsedPlugin, error) {
	var plugins []*ParsedPlugin

	// Walk through repository looking for manifest.json files
	err := filepath.WalkDir(repoPath, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		// Skip .git directory
		if d.IsDir() && d.Name() == ".git" {
			return filepath.SkipDir
		}

		// Only process manifest.json files
		if !d.IsDir() && d.Name() == "manifest.json" {
			plugin, err := p.ParsePluginFile(path)
			if err != nil {
				// Log error but continue processing other files
				return nil
			}

			plugins = append(plugins, plugin)
		}

		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to walk repository: %w", err)
	}

	return plugins, nil
}

// ParsePluginFile parses a single plugin manifest.json file
func (p *PluginParser) ParsePluginFile(filePath string) (*ParsedPlugin, error) {
	// Read file
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	// Parse JSON
	var manifest PluginManifest
	if err := json.Unmarshal(data, &manifest); err != nil {
		return nil, fmt.Errorf("failed to parse JSON: %w", err)
	}

	// Validate required fields
	if manifest.Name == "" {
		return nil, fmt.Errorf("plugin name is required")
	}

	if manifest.Version == "" {
		return nil, fmt.Errorf("version is required")
	}

	if manifest.DisplayName == "" {
		return nil, fmt.Errorf("displayName is required")
	}

	if manifest.Type == "" {
		return nil, fmt.Errorf("type is required")
	}

	// Validate plugin type
	validTypes := map[string]bool{
		"extension": true,
		"webhook":   true,
		"api":       true,
		"ui":        true,
		"theme":     true,
	}
	if !validTypes[manifest.Type] {
		return nil, fmt.Errorf("invalid plugin type: %s (must be extension, webhook, api, ui, or theme)", manifest.Type)
	}

	// Convert full manifest to JSON for storage
	manifestJSON, err := json.Marshal(manifest)
	if err != nil {
		return nil, fmt.Errorf("failed to encode manifest: %w", err)
	}

	plugin := &ParsedPlugin{
		Name:        manifest.Name,
		Version:     manifest.Version,
		DisplayName: manifest.DisplayName,
		Description: manifest.Description,
		Category:    manifest.Category,
		PluginType:  manifest.Type,
		Icon:        manifest.Icon,
		Manifest:    string(manifestJSON),
		Tags:        manifest.Tags,
	}

	if plugin.Tags == nil {
		plugin.Tags = []string{}
	}

	return plugin, nil
}

// ValidatePluginManifest validates a plugin manifest structure
func (p *PluginParser) ValidatePluginManifest(jsonContent string) error {
	var manifest PluginManifest
	if err := json.Unmarshal([]byte(jsonContent), &manifest); err != nil {
		return fmt.Errorf("invalid JSON: %w", err)
	}

	// Check required fields
	if manifest.Name == "" {
		return fmt.Errorf("name is required")
	}

	if manifest.Version == "" {
		return fmt.Errorf("version is required")
	}

	if manifest.DisplayName == "" {
		return fmt.Errorf("displayName is required")
	}

	if manifest.Type == "" {
		return fmt.Errorf("type is required")
	}

	// Validate plugin type
	validTypes := map[string]bool{
		"extension": true,
		"webhook":   true,
		"api":       true,
		"ui":        true,
		"theme":     true,
	}
	if !validTypes[manifest.Type] {
		return fmt.Errorf("invalid type: %s (must be extension, webhook, api, ui, or theme)", manifest.Type)
	}

	return nil
}
