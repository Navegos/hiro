// Copyright 2023 Heroic Labs & Contributors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package hiro

import (
	"context"
	"encoding/json"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/heroiclabs/nakama-common/runtime"
)

type SatoriPublisher interface {
	IsPublishAuthenticateRequest() bool
	IsPublishAchievementsEvents() bool
	IsPublishBaseEvents() bool
	IsPublishEconomyEvents() bool
	IsPublishEnergyEvents() bool
	IsPublishEventLeaderboardsEvents() bool
	IsPublishIncentivesEvents() bool
	IsPublishInventoryEvents() bool
	IsPublishLeaderboardsEvents() bool
	IsPublishProgressionEvents() bool
	IsPublishStatsEvents() bool
	IsPublishTeamsEvents() bool
	IsPublishTutorialsEvents() bool
	IsPublishUnlockablesEvents() bool
}

var _ SatoriPublisher = (*SatoriPersonalizer)(nil)

var _ Personalizer = (*SatoriPersonalizer)(nil)

type SatoriPersonalizerOption interface {
	apply(*SatoriPersonalizer)
}

type satoriPersonalizerOptionFunc struct {
	f func(*SatoriPersonalizer)
}

func (s *satoriPersonalizerOptionFunc) apply(personalizer *SatoriPersonalizer) {
	s.f(personalizer)
}

func SatoriPersonalizerPublishAuthenticateEvents() SatoriPersonalizerOption {
	return &satoriPersonalizerOptionFunc{
		f: func(personalizer *SatoriPersonalizer) {
			personalizer.publishAuthenticateRequest = true
		},
	}
}

func SatoriPersonalizerPublishAchievementsEvents() SatoriPersonalizerOption {
	return &satoriPersonalizerOptionFunc{
		f: func(personalizer *SatoriPersonalizer) {
			personalizer.publishAchievementsEvents = true
		},
	}
}

func SatoriPersonalizerPublishBaseEvents() SatoriPersonalizerOption {
	return &satoriPersonalizerOptionFunc{
		f: func(personalizer *SatoriPersonalizer) {
			personalizer.publishBaseEvents = true
		},
	}
}

func SatoriPersonalizerPublishEconomyEvents() SatoriPersonalizerOption {
	return &satoriPersonalizerOptionFunc{
		f: func(personalizer *SatoriPersonalizer) {
			personalizer.publishEconomyEvents = true
		},
	}
}

func SatoriPersonalizerPublishEnergyEvents() SatoriPersonalizerOption {
	return &satoriPersonalizerOptionFunc{
		f: func(personalizer *SatoriPersonalizer) {
			personalizer.publishEnergyEvents = true
		},
	}
}

func SatoriPersonalizerPublishEventLeaderboardsEvents() SatoriPersonalizerOption {
	return &satoriPersonalizerOptionFunc{
		f: func(personalizer *SatoriPersonalizer) {
			personalizer.publishEventLeaderboardsEvents = true
		},
	}
}

func SatoriPersonalizerPublishIncentivesEvents() SatoriPersonalizerOption {
	return &satoriPersonalizerOptionFunc{
		f: func(personalizer *SatoriPersonalizer) {
			personalizer.publishIncentivesEvents = true
		},
	}
}

func SatoriPersonalizerPublishInventoryEvents() SatoriPersonalizerOption {
	return &satoriPersonalizerOptionFunc{
		f: func(personalizer *SatoriPersonalizer) {
			personalizer.publishInventoryEvents = true
		},
	}
}

func SatoriPersonalizerPublishLeaderboardsEvents() SatoriPersonalizerOption {
	return &satoriPersonalizerOptionFunc{
		f: func(personalizer *SatoriPersonalizer) {
			personalizer.publishLeaderboardsEvents = true
		},
	}
}

func SatoriPersonalizerPublishProgressionEvents() SatoriPersonalizerOption {
	return &satoriPersonalizerOptionFunc{
		f: func(personalizer *SatoriPersonalizer) {
			personalizer.publishProgressionEvents = true
		},
	}
}

func SatoriPersonalizerPublishStatsEvents() SatoriPersonalizerOption {
	return &satoriPersonalizerOptionFunc{
		f: func(personalizer *SatoriPersonalizer) {
			personalizer.publishStatsEvents = true
		},
	}
}

func SatoriPersonalizerPublishTeamsEvents() SatoriPersonalizerOption {
	return &satoriPersonalizerOptionFunc{
		f: func(personalizer *SatoriPersonalizer) {
			personalizer.publishTeamsEvents = true
		},
	}
}

func SatoriPersonalizerPublishTutorialsEvents() SatoriPersonalizerOption {
	return &satoriPersonalizerOptionFunc{
		f: func(personalizer *SatoriPersonalizer) {
			personalizer.publishTutorialsEvents = true
		},
	}
}

func SatoriPersonalizerPublishUnlockablesEvents() SatoriPersonalizerOption {
	return &satoriPersonalizerOptionFunc{
		f: func(personalizer *SatoriPersonalizer) {
			personalizer.publishUnlockablesEvents = true
		},
	}
}

func SatoriPersonalizerNoCache() SatoriPersonalizerOption {
	return &satoriPersonalizerOptionFunc{
		f: func(personalizer *SatoriPersonalizer) {
			personalizer.noCache = true
		},
	}
}

type SatoriPersonalizerCache struct {
	flags      *runtime.FlagList
	liveEvents *atomic.Pointer[runtime.LiveEventList]
}

type SatoriPersonalizer struct {
	publishAuthenticateRequest bool

	publishAchievementsEvents      bool
	publishBaseEvents              bool
	publishEconomyEvents           bool
	publishEnergyEvents            bool
	publishEventLeaderboardsEvents bool
	publishIncentivesEvents        bool
	publishInventoryEvents         bool
	publishLeaderboardsEvents      bool
	publishProgressionEvents       bool
	publishStatsEvents             bool
	publishTeamsEvents             bool
	publishTutorialsEvents         bool
	publishUnlockablesEvents       bool

	noCache bool

	cacheMutex sync.RWMutex
	cache      map[context.Context]*SatoriPersonalizerCache
}

func NewSatoriPersonalizer(ctx context.Context, opts ...SatoriPersonalizerOption) *SatoriPersonalizer {
	s := &SatoriPersonalizer{
		cacheMutex: sync.RWMutex{},
		cache:      make(map[context.Context]*SatoriPersonalizerCache),
	}

	// Apply options, if any supplied.
	for _, opt := range opts {
		opt.apply(s)
	}

	if !s.noCache {
		go func() {
			ticker := time.NewTicker(30 * time.Second)
			defer ticker.Stop()
			for {
				select {
				case <-ctx.Done():
					return
				case <-ticker.C:
					s.cacheMutex.Lock()
					for cacheCtx := range s.cache {
						if cacheCtx.Err() != nil {
							delete(s.cache, cacheCtx)
						}
					}
					s.cacheMutex.Unlock()
				}
			}
		}()
	}

	return s
}

var allFlagNames = []string{"Hiro-Achievements", "Hiro-Base", "Hiro-Economy", "Hiro-Energy", "Hiro-Inventory", "Hiro-Leaderboards", "Hiro-Teams", "Hiro-Tutorials", "Hiro-Unlockables", "Hiro-Stats", "Hiro-Event-Leaderboards", "Hiro-Progression", "Hiro-Incentives"}

func (p *SatoriPersonalizer) GetValue(ctx context.Context, logger runtime.Logger, nk runtime.NakamaModule, system System, userID string) (any, error) {
	var flagName string
	switch system.GetType() {
	case SystemTypeAchievements:
		flagName = "Hiro-Achievements"
	case SystemTypeBase:
		flagName = "Hiro-Base"
	case SystemTypeEconomy:
		flagName = "Hiro-Economy"
	case SystemTypeEnergy:
		flagName = "Hiro-Energy"
	case SystemTypeInventory:
		flagName = "Hiro-Inventory"
	case SystemTypeLeaderboards:
		flagName = "Hiro-Leaderboards"
	case SystemTypeTeams:
		flagName = "Hiro-Teams"
	case SystemTypeTutorials:
		flagName = "Hiro-Tutorials"
	case SystemTypeUnlockables:
		flagName = "Hiro-Unlockables"
	case SystemTypeStats:
		flagName = "Hiro-Stats"
	case SystemTypeEventLeaderboards:
		flagName = "Hiro-Event-Leaderboards"
	case SystemTypeProgression:
		flagName = "Hiro-Progression"
	case SystemTypeIncentives:
		flagName = "Hiro-Incentives"
	default:
		return nil, runtime.NewError("hiro system type unknown", 3)
	}

	var config any
	var found bool

	if p.noCache {
		flagList, err := nk.GetSatori().FlagsList(ctx, userID, flagName)
		if err != nil {
			if strings.Contains(err.Error(), "404 status code") {
				logger.WithField("userID", userID).WithField("error", err.Error()).Warn("error requesting Satori flag list, user not found")
				return nil, nil
			}
			logger.WithField("userID", userID).WithField("error", err.Error()).Error("error requesting Satori flag list")
			return nil, err
		}

		if len(flagList.Flags) >= 1 {
			config = system.GetConfig()
			decoder := json.NewDecoder(strings.NewReader(flagList.Flags[0].Value))
			decoder.DisallowUnknownFields()
			if err := decoder.Decode(config); err != nil {
				logger.WithField("userID", userID).WithField("error", err.Error()).Error("error merging Satori flag value")
				return nil, err
			}
			found = true
		}

		if s := system.GetType(); s == SystemTypeEventLeaderboards || s == SystemTypeAchievements {
			// If looking at event leaderboards, also load live events.
			liveEventsList, err := nk.GetSatori().LiveEventsList(ctx, userID)
			if err != nil {
				if strings.Contains(err.Error(), "404 status code") {
					logger.WithField("userID", userID).WithField("error", err.Error()).Warn("error requesting Satori live events list, user not found")
					return nil, nil
				}
				logger.WithField("userID", userID).WithField("error", err.Error()).Error("error requesting Satori live events list")
				return nil, err
			}
			if len(liveEventsList.LiveEvents) > 0 {
				if config == nil {
					config = system.GetConfig()
				}
				for _, liveEvent := range liveEventsList.LiveEvents {
					decoder := json.NewDecoder(strings.NewReader(liveEvent.Value))
					decoder.DisallowUnknownFields()
					if err := decoder.Decode(config); err != nil {
						// The live event may be intended for a different purpose, do not log or return an error here.
						continue
					}
					found = true
				}
			}
		}
	} else {
		var cacheEntry *SatoriPersonalizerCache
		p.cacheMutex.RLock()
		cacheEntry, found = p.cache[ctx]
		p.cacheMutex.RUnlock()

		if !found {
			flagList, err := nk.GetSatori().FlagsList(ctx, userID, allFlagNames...)
			if err != nil {
				if strings.Contains(err.Error(), "404 status code") {
					logger.WithField("userID", userID).WithField("error", err.Error()).Warn("error requesting Satori flag list, user not found")
					return nil, nil
				}
				logger.WithField("userID", userID).WithField("error", err.Error()).Error("error requesting Satori flag list")
				return nil, err
			}

			var liveEventsList *runtime.LiveEventList
			if s := system.GetType(); s == SystemTypeEventLeaderboards || s == SystemTypeAchievements {
				liveEventsList, err = nk.GetSatori().LiveEventsList(ctx, userID)
				if err != nil {
					if strings.Contains(err.Error(), "404 status code") {
						logger.WithField("userID", userID).WithField("error", err.Error()).Warn("error requesting Satori live events list, user not found")
						return nil, nil
					}
					logger.WithField("userID", userID).WithField("error", err.Error()).Error("error requesting Satori live events list")
					return nil, err
				}
			}

			cacheEntry = &SatoriPersonalizerCache{
				flags:      flagList,
				liveEvents: &atomic.Pointer[runtime.LiveEventList]{},
			}
			if liveEventsList != nil {
				cacheEntry.liveEvents.Store(liveEventsList)
			}
			p.cacheMutex.Lock()
			p.cache[ctx] = cacheEntry
			p.cacheMutex.Unlock()
		}

		if s := system.GetType(); (s == SystemTypeEventLeaderboards || s == SystemTypeAchievements) && cacheEntry.liveEvents.Load() == nil {
			liveEventsList, err := nk.GetSatori().LiveEventsList(ctx, userID)
			if err != nil {
				if strings.Contains(err.Error(), "404 status code") {
					logger.WithField("userID", userID).WithField("error", err.Error()).Warn("error requesting Satori live events list, user not found")
					return nil, nil
				}
				logger.WithField("userID", userID).WithField("error", err.Error()).Error("error requesting Satori live events list")
				return nil, err
			}
			cacheEntry.liveEvents.Store(liveEventsList)
		}

		found = false

		for _, flag := range cacheEntry.flags.Flags {
			if flag.Name != flagName {
				continue
			}

			config = system.GetConfig()
			decoder := json.NewDecoder(strings.NewReader(flag.Value))
			decoder.DisallowUnknownFields()
			if err := decoder.Decode(config); err != nil {
				logger.WithField("userID", userID).WithField("error", err.Error()).Error("error merging Satori flag value")
				return nil, err
			}
			found = true
		}

		if liveEventsList := cacheEntry.liveEvents.Load(); liveEventsList != nil && len(liveEventsList.LiveEvents) > 0 {
			if config == nil {
				config = system.GetConfig()
			}
			for _, liveEvent := range liveEventsList.LiveEvents {
				decoder := json.NewDecoder(strings.NewReader(liveEvent.Value))
				decoder.DisallowUnknownFields()
				if err := decoder.Decode(config); err != nil {
					// The live event may be intended for a different purpose, do not log or return an error here.
					continue
				}
				found = true
			}
		}
	}

	// If this caller doesn't have the given flag (or live events) return the nil to indicate no change to the config.
	if !found {
		return nil, nil
	}

	return config, nil
}

func (p *SatoriPersonalizer) IsPublishAuthenticateRequest() bool {
	return p.publishAuthenticateRequest
}

func (p *SatoriPersonalizer) IsPublishAchievementsEvents() bool {
	return p.publishAchievementsEvents
}

func (p *SatoriPersonalizer) IsPublishBaseEvents() bool {
	return p.publishBaseEvents
}

func (p *SatoriPersonalizer) IsPublishEconomyEvents() bool {
	return p.publishEconomyEvents
}

func (p *SatoriPersonalizer) IsPublishEnergyEvents() bool {
	return p.publishEnergyEvents
}

func (p *SatoriPersonalizer) IsPublishEventLeaderboardsEvents() bool {
	return p.publishEventLeaderboardsEvents
}

func (p *SatoriPersonalizer) IsPublishIncentivesEvents() bool {
	return p.publishIncentivesEvents
}

func (p *SatoriPersonalizer) IsPublishInventoryEvents() bool {
	return p.publishInventoryEvents
}

func (p *SatoriPersonalizer) IsPublishLeaderboardsEvents() bool {
	return p.publishLeaderboardsEvents
}

func (p *SatoriPersonalizer) IsPublishProgressionEvents() bool {
	return p.publishProgressionEvents
}

func (p *SatoriPersonalizer) IsPublishStatsEvents() bool {
	return p.publishStatsEvents
}

func (p *SatoriPersonalizer) IsPublishTeamsEvents() bool {
	return p.publishTeamsEvents
}

func (p *SatoriPersonalizer) IsPublishTutorialsEvents() bool {
	return p.publishTutorialsEvents
}

func (p *SatoriPersonalizer) IsPublishUnlockablesEvents() bool {
	return p.publishUnlockablesEvents
}
