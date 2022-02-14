package recipes

import (
	"context"

	"github.com/newrelic/newrelic-cli/internal/install/execution"
	"github.com/newrelic/newrelic-cli/internal/install/types"
)

type DetectionStatusProvider interface {
	DetectionStatus(context.Context, *types.OpenInstallationRecipe) execution.RecipeStatusType
}

type RecipeDetector struct {
	processEvaluator DetectionStatusProvider
	scriptEvaluator  DetectionStatusProvider
}

func newRecipeDetector(processEvaluator DetectionStatusProvider, scriptEvaluator DetectionStatusProvider) *RecipeDetector {
	return &RecipeDetector{
		processEvaluator: processEvaluator,
		scriptEvaluator:  scriptEvaluator,
	}
}

func NewRecipeDetector() *RecipeDetector {
	return newRecipeDetector(NewProcessEvaluator(), NewScriptEvaluator())
}

func (dt *RecipeDetector) DetectBundle(ctx context.Context, bundle *Bundle) *Bundle {

	for _, bundleRecipe := range bundle.BundleRecipes {
		status := dt.detectRecipe(ctx, bundleRecipe.Recipe)
		bundleRecipe.Statuses = append(bundleRecipe.Statuses, status)
	}

	return bundle
}

//TODO: might not need?
func (dt *RecipeDetector) DetectRecipes(ctx context.Context, recipes []*types.OpenInstallationRecipe) map[*types.OpenInstallationRecipe]execution.RecipeStatusType {

	results := make(map[*types.OpenInstallationRecipe]execution.RecipeStatusType)

	for _, recipe := range recipes {
		results[recipe] = dt.detectRecipe(ctx, recipe)
	}

	return results
}

func (dt *RecipeDetector) detectRecipe(ctx context.Context, recipe *types.OpenInstallationRecipe) execution.RecipeStatusType {

	status := dt.processEvaluator.DetectionStatus(ctx, recipe)

	if status == execution.RecipeStatusTypes.AVAILABLE && recipe.PreInstall.RequireAtDiscovery != "" {
		status = dt.scriptEvaluator.DetectionStatus(ctx, recipe)
	}

	return status
}
