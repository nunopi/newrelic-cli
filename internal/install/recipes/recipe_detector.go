package recipes

import (
	"context"
	"time"

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

//TODO: Should this also detect recipe dependencies?  Should recipe itself be available to install if one of its dependency is not?
func (dt *RecipeDetector) DetectBundle(ctx context.Context, bundle *Bundle) {
	for _, bundleRecipe := range bundle.BundleRecipes {
		status := dt.detectRecipe(ctx, bundleRecipe.Recipe)
		bundleRecipe.AddStatus(status, time.Now())
	}
}

func (dt *RecipeDetector) detectRecipe(ctx context.Context, recipe *types.OpenInstallationRecipe) execution.RecipeStatusType {

	status := dt.processEvaluator.DetectionStatus(ctx, recipe)

	if status == execution.RecipeStatusTypes.AVAILABLE && recipe.PreInstall.RequireAtDiscovery != "" {
		status = dt.scriptEvaluator.DetectionStatus(ctx, recipe)
	}

	return status
}