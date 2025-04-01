package gocritic

import (
	"errors"
	"fmt"
	"maps"
	"slices"
	"strings"

	gocriticlinter "github.com/go-critic/go-critic/linter"

	"github.com/golangci/golangci-lint/v2/pkg/config"
	"github.com/golangci/golangci-lint/v2/pkg/logutils"
)

type settingsWrapper struct {
	*config.GoCriticSettings

	logger logutils.Log

	allCheckers []*gocriticlinter.CheckerInfo

	allChecks             goCriticChecks[struct{}]
	allChecksByTag        goCriticChecks[[]string]
	allTagsSorted         []string
	inferredEnabledChecks goCriticChecks[struct{}]

	// *LowerCased fields are used for GoCriticSettings.SettingsPerCheck validation only.

	allChecksLowerCased             goCriticChecks[struct{}]
	inferredEnabledChecksLowerCased goCriticChecks[struct{}]
}

func newSettingsWrapper(settings *config.GoCriticSettings, logger logutils.Log) *settingsWrapper {
	allCheckers := gocriticlinter.GetCheckersInfo()

	allChecks := make(goCriticChecks[struct{}], len(allCheckers))
	allChecksLowerCased := make(goCriticChecks[struct{}], len(allCheckers))
	allChecksByTag := make(goCriticChecks[[]string])
	for _, checker := range allCheckers {
		allChecks[checker.Name] = struct{}{}
		allChecksLowerCased[strings.ToLower(checker.Name)] = struct{}{}

		for _, tag := range checker.Tags {
			allChecksByTag[tag] = append(allChecksByTag[tag], checker.Name)
		}
	}

	allTagsSorted := slices.Sorted(maps.Keys(allChecksByTag))

	return &settingsWrapper{
		GoCriticSettings:                settings,
		logger:                          logger,
		allCheckers:                     allCheckers,
		allChecks:                       allChecks,
		allChecksLowerCased:             allChecksLowerCased,
		allChecksByTag:                  allChecksByTag,
		allTagsSorted:                   allTagsSorted,
		inferredEnabledChecks:           make(goCriticChecks[struct{}]),
		inferredEnabledChecksLowerCased: make(goCriticChecks[struct{}]),
	}
}

func (s *settingsWrapper) IsCheckEnabled(name string) bool {
	return s.inferredEnabledChecks.has(name)
}

func (s *settingsWrapper) GetLowerCasedParams() map[string]config.GoCriticCheckSettings {
	return normalizeMap(s.SettingsPerCheck)
}

// InferEnabledChecks tries to be consistent with (lintersdb.Manager).build.
func (s *settingsWrapper) InferEnabledChecks() {
	s.debugChecksInitialState()

	enabledByDefaultChecks, disabledByDefaultChecks := s.buildEnabledAndDisabledByDefaultChecks()

	debugChecksListf(enabledByDefaultChecks, "Enabled by default")
	debugChecksListf(disabledByDefaultChecks, "Disabled by default")

	enabledChecks := make(goCriticChecks[struct{}])

	if s.EnableAll {
		enabledChecks = make(goCriticChecks[struct{}], len(s.allCheckers))
		for _, info := range s.allCheckers {
			enabledChecks[info.Name] = struct{}{}
		}
	} else if !s.DisableAll {
		// enable-all/disable-all revokes the default settings.
		enabledChecks = make(goCriticChecks[struct{}], len(enabledByDefaultChecks))
		for _, check := range enabledByDefaultChecks {
			enabledChecks[check] = struct{}{}
		}
	}

	if len(s.EnabledTags) != 0 {
		enabledFromTags := s.expandTagsToChecks(s.EnabledTags)

		debugChecksListf(enabledFromTags, "Enabled by config tags %s", s.EnabledTags)

		for _, check := range enabledFromTags {
			enabledChecks[check] = struct{}{}
		}
	}

	if len(s.EnabledChecks) != 0 {
		debugChecksListf(s.EnabledChecks, "Enabled by config")

		for _, check := range s.EnabledChecks {
			if enabledChecks.has(check) {
				s.logger.Warnf("%s: no need to enable check %q: it's already enabled", linterName, check)
				continue
			}
			enabledChecks[check] = struct{}{}
		}
	}

	if len(s.DisabledTags) != 0 {
		disabledFromTags := s.expandTagsToChecks(s.DisabledTags)

		debugChecksListf(disabledFromTags, "Disabled by config tags %s", s.DisabledTags)

		for _, check := range disabledFromTags {
			delete(enabledChecks, check)
		}
	}

	if len(s.DisabledChecks) != 0 {
		debugChecksListf(s.DisabledChecks, "Disabled by config")

		for _, check := range s.DisabledChecks {
			if !enabledChecks.has(check) {
				s.logger.Warnf("%s: no need to disable check %q: it's already disabled", linterName, check)
				continue
			}
			delete(enabledChecks, check)
		}
	}

	s.inferredEnabledChecks = enabledChecks
	s.inferredEnabledChecksLowerCased = normalizeMap(s.inferredEnabledChecks)

	s.debugChecksFinalState()
}

func (s *settingsWrapper) buildEnabledAndDisabledByDefaultChecks() (enabled, disabled []string) {
	for _, info := range s.allCheckers {
		if enabledByDef := isEnabledByDefaultGoCriticChecker(info); enabledByDef {
			enabled = append(enabled, info.Name)
		} else {
			disabled = append(disabled, info.Name)
		}
	}
	return enabled, disabled
}

func (s *settingsWrapper) expandTagsToChecks(tags []string) []string {
	var checks []string
	for _, tag := range tags {
		checks = append(checks, s.allChecksByTag[tag]...)
	}
	return checks
}

func (s *settingsWrapper) debugChecksInitialState() {
	if !isDebug {
		return
	}

	debugf("All gocritic existing tags and checks:")
	for _, tag := range s.allTagsSorted {
		debugChecksListf(s.allChecksByTag[tag], "  tag %q", tag)
	}
}

func (s *settingsWrapper) debugChecksFinalState() {
	if !isDebug {
		return
	}

	var enabledChecks []string
	var disabledChecks []string

	for _, checker := range s.allCheckers {
		check := checker.Name
		if s.inferredEnabledChecks.has(check) {
			enabledChecks = append(enabledChecks, check)
		} else {
			disabledChecks = append(disabledChecks, check)
		}
	}

	debugChecksListf(enabledChecks, "Final used")

	if len(disabledChecks) == 0 {
		debugf("All checks are enabled")
	} else {
		debugChecksListf(disabledChecks, "Final not used")
	}
}

// Validate tries to be consistent with (lintersdb.Validator).validateEnabledDisabledLintersConfig.
func (s *settingsWrapper) Validate() error {
	for _, v := range []func() error{
		s.validateOptionsCombinations,
		s.validateCheckerTags,
		s.validateCheckerNames,
		s.validateDisabledAndEnabledAtOneMoment,
		s.validateAtLeastOneCheckerEnabled,
	} {
		if err := v(); err != nil {
			return err
		}
	}
	return nil
}

func (s *settingsWrapper) validateOptionsCombinations() error {
	if s.EnableAll {
		if s.DisableAll {
			return errors.New("enable-all and disable-all options must not be combined")
		}

		if len(s.EnabledTags) != 0 {
			return errors.New("enable-all and enabled-tags options must not be combined")
		}

		if len(s.EnabledChecks) != 0 {
			return errors.New("enable-all and enabled-checks options must not be combined")
		}
	}

	if s.DisableAll {
		if len(s.DisabledTags) != 0 {
			return errors.New("disable-all and disabled-tags options must not be combined")
		}

		if len(s.DisabledChecks) != 0 {
			return errors.New("disable-all and disabled-checks options must not be combined")
		}

		if len(s.EnabledTags) == 0 && len(s.EnabledChecks) == 0 {
			return errors.New("all checks were disabled, but no one check was enabled: at least one must be enabled")
		}
	}

	return nil
}

func (s *settingsWrapper) validateCheckerTags() error {
	for _, tag := range s.EnabledTags {
		if !s.allChecksByTag.has(tag) {
			return fmt.Errorf("enabled tag %q doesn't exist, see %s's documentation", tag, linterName)
		}
	}

	for _, tag := range s.DisabledTags {
		if !s.allChecksByTag.has(tag) {
			return fmt.Errorf("disabled tag %q doesn't exist, see %s's documentation", tag, linterName)
		}
	}

	return nil
}

func (s *settingsWrapper) validateCheckerNames() error {
	for _, check := range s.EnabledChecks {
		if !s.allChecks.has(check) {
			return fmt.Errorf("enabled check %q doesn't exist, see %s's documentation", check, linterName)
		}
	}

	for _, check := range s.DisabledChecks {
		if !s.allChecks.has(check) {
			return fmt.Errorf("disabled check %q doesn't exist, see %s documentation", check, linterName)
		}
	}

	for check := range s.SettingsPerCheck {
		lcName := strings.ToLower(check)
		if !s.allChecksLowerCased.has(lcName) {
			return fmt.Errorf("invalid check settings: check %q doesn't exist, see %s documentation", check, linterName)
		}
		if !s.inferredEnabledChecksLowerCased.has(lcName) {
			s.logger.Warnf("%s: settings were provided for disabled check %q", check, linterName)
		}
	}

	return nil
}

func (s *settingsWrapper) validateDisabledAndEnabledAtOneMoment() error {
	for _, tag := range s.DisabledTags {
		if slices.Contains(s.EnabledTags, tag) {
			return fmt.Errorf("tag %q disabled and enabled at one moment", tag)
		}
	}

	for _, check := range s.DisabledChecks {
		if slices.Contains(s.EnabledChecks, check) {
			return fmt.Errorf("check %q disabled and enabled at one moment", check)
		}
	}

	return nil
}

func (s *settingsWrapper) validateAtLeastOneCheckerEnabled() error {
	if len(s.inferredEnabledChecks) == 0 {
		return errors.New("eventually all checks were disabled: at least one must be enabled")
	}
	return nil
}

type goCriticChecks[T any] map[string]T

func (m goCriticChecks[T]) has(name string) bool {
	_, ok := m[name]
	return ok
}

func debugChecksListf(checks []string, format string, args ...any) {
	if !isDebug {
		return
	}

	v := slices.Sorted(slices.Values(checks))

	debugf("%s checks (%d): %s", fmt.Sprintf(format, args...), len(checks), strings.Join(v, ", "))
}
