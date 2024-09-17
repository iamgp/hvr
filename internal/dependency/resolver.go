package dependency

import (
	"fmt"

	"github.com/Masterminds/semver/v3"
	"github.com/iamgp/hvr/internal/models"
	"github.com/iamgp/hvr/internal/storage"
)

type Resolver struct {
	db *storage.SQLiteDatabase
}

func NewResolver(db *storage.SQLiteDatabase) *Resolver {
	return &Resolver{db: db}
}

func (r *Resolver) ResolveDependencies(library models.Library) ([]models.Library, error) {
	resolved := make(map[string]models.Library)
	err := r.resolveDependenciesRecursive(library, resolved, 0)
	if err != nil {
		return nil, err
	}

	result := make([]models.Library, 0, len(resolved))
	for _, lib := range resolved {
		result = append(result, lib)
	}
	return result, nil
}

func (r *Resolver) resolveDependenciesRecursive(library models.Library, resolved map[string]models.Library, depth int) error {
	if depth > 100 {
		return fmt.Errorf("dependency resolution too deep, possible circular dependency")
	}

	for depName, depVersionConstraint := range library.Dependencies {
		constraint, err := semver.NewConstraint(depVersionConstraint)
		if err != nil {
			return fmt.Errorf("invalid version constraint for %s: %w", depName, err)
		}

		versions, err := r.db.GetAllVersions(depName)
		if err != nil {
			return fmt.Errorf("failed to get versions for %s: %w", depName, err)
		}

		var bestMatch *semver.Version
		for _, version := range versions {
			if constraint.Check(version) {
				if bestMatch == nil || version.GreaterThan(bestMatch) {
					bestMatch = version
				}
			}
		}

		if bestMatch == nil {
			return fmt.Errorf("no suitable version found for %s matching %s", depName, depVersionConstraint)
		}

		depLibrary, err := r.db.Get(depName, bestMatch.String())
		if err != nil {
			return fmt.Errorf("failed to get library %s version %s: %w", depName, bestMatch.String(), err)
		}

		if _, exists := resolved[depName]; !exists {
			resolved[depName] = depLibrary
			err = r.resolveDependenciesRecursive(depLibrary, resolved, depth+1)
			if err != nil {
				return err
			}
		}
	}

	return nil
}
