package servermanager

import (
	"math/rand"
	"sort"
	"time"

	"github.com/google/uuid"
)

// RaceWeekendEntryListSorter is a function which takes a race weekend, session and entrylist and sorts the entrylist based on some criteria.
<<<<<<< HEAD
type RaceWeekendEntryListSorter interface {
	Sort(*RaceWeekend, *RaceWeekendSession, []*RaceWeekendSessionEntrant, *RaceWeekendSessionToSessionFilter) error
}

type RaceWeekendEntryListSorterDescription struct {
	Name                  string
	Key                   string
	Sorter                RaceWeekendEntryListSorter
	NeedsParentSession    bool
	NeedsChampionship     bool
	ShowInManageEntryList bool
}

type RaceWeekendEntryListSortFunc func(*RaceWeekend, *RaceWeekendSession, []*RaceWeekendSessionEntrant, *RaceWeekendSessionToSessionFilter) error

func (rwelsf RaceWeekendEntryListSortFunc) Sort(rw *RaceWeekend, rws *RaceWeekendSession, rwes []*RaceWeekendSessionEntrant, rwsf *RaceWeekendSessionToSessionFilter) error {
	return rwelsf(rw, rws, rwes, rwsf)
=======
type RaceWeekendEntryListSorter func(*RaceWeekendSession, []*RaceWeekendSessionEntrant) error

type RaceWeekendEntryListSorterDescription struct {
	Name     string
	Key      string
	SortFunc RaceWeekendEntryListSorter
>>>>>>> origin/multiserver2
}

var RaceWeekendEntryListSorters = []RaceWeekendEntryListSorterDescription{
	{
<<<<<<< HEAD
		Name:                  "No Sort (Use Finishing Grid)",
		Key:                   "", // key intentionally left blank
		Sorter:                RaceWeekendEntryListSortFunc(UnchangedRaceWeekendEntryListSort),
		NeedsParentSession:    false,
		NeedsChampionship:     false,
		ShowInManageEntryList: true,
	},
	{
		Name:                  "Fastest Lap",
		Key:                   "fastest_lap",
		Sorter:                RaceWeekendEntryListSortFunc(FastestLapRaceWeekendEntryListSort),
		NeedsParentSession:    true,
		NeedsChampionship:     false,
		ShowInManageEntryList: true,
	},
	{
		Name:                  "Total Race Time",
		Key:                   "total_race_time",
		Sorter:                RaceWeekendEntryListSortFunc(TotalRaceTimeRaceWeekendEntryListSort),
		NeedsParentSession:    true,
		NeedsChampionship:     false,
		ShowInManageEntryList: true,
	},
	{
		Name:                  "Fastest Lap Across Multiple Results Files",
		Key:                   "fastest_multi_results_lap",
		Sorter:                RaceWeekendEntryListSortFunc(FastestResultsFileRaceWeekendEntryListSort),
		NeedsParentSession:    false,
		NeedsChampionship:     false,
		ShowInManageEntryList: false,
	},
	{
		Name:                  "Number of Laps Across Multiple Results Files",
		Key:                   "number_multi_results_lap",
		Sorter:                RaceWeekendEntryListSortFunc(NumberResultsFileRaceWeekendEntryListSort),
		NeedsParentSession:    false,
		NeedsChampionship:     false,
		ShowInManageEntryList: false,
	},
	{
		Name:                  "Fewest Collisions",
		Key:                   "fewest_collisions",
		Sorter:                RaceWeekendEntryListSortFunc(FewestCollisionsRaceWeekendEntryListSort),
		NeedsParentSession:    true,
		NeedsChampionship:     false,
		ShowInManageEntryList: true,
	},
	{
		Name:                  "Fewest Cuts",
		Key:                   "fewest_cuts",
		Sorter:                RaceWeekendEntryListSortFunc(FewestCutsRaceWeekendEntryListSort),
		NeedsParentSession:    true,
		NeedsChampionship:     false,
		ShowInManageEntryList: true,
	},
	{
		Name:                  "Safety (Collisions then Cuts)",
		Key:                   "safety",
		Sorter:                RaceWeekendEntryListSortFunc(SafetyRaceWeekendEntryListSort),
		NeedsParentSession:    true,
		NeedsChampionship:     false,
		ShowInManageEntryList: true,
	},
	{
		Name:                  "Championship Standings Order",
		Key:                   "championship_standings_order",
		Sorter:                &ChampionshipStandingsOrderEntryListSort{},
		NeedsParentSession:    false,
		NeedsChampionship:     true,
		ShowInManageEntryList: true,
	},
	{
		Name:                  "Championship Class",
		Key:                   "championship_class",
		Sorter:                &ChampionshipClassSort{},
		NeedsParentSession:    false,
		NeedsChampionship:     true,
		ShowInManageEntryList: false,
	},
	{
		Name:                  "Random",
		Key:                   "random",
		Sorter:                RaceWeekendEntryListSortFunc(RandomRaceWeekendEntryListSort),
		NeedsParentSession:    false,
		NeedsChampionship:     false,
		ShowInManageEntryList: true,
	},
	{
		Name:                  "Alphabetical (Using Driver Name)",
		Key:                   "alphabetical",
		Sorter:                RaceWeekendEntryListSortFunc(AlphabeticalRaceWeekendEntryListSort),
		NeedsParentSession:    false,
		NeedsChampionship:     false,
		ShowInManageEntryList: true,
=======
		Name:     "No Sort (Use Finishing Grid)",
		Key:      "", // key intentionally left blank
		SortFunc: UnchangedRaceWeekendEntryListSort,
	},
	{
		Name:     "Fastest Lap",
		Key:      "fastest_lap",
		SortFunc: FastestLapRaceWeekendEntryListSort,
	},
	{
		Name:     "Total Race Time",
		Key:      "total_race_time",
		SortFunc: TotalRaceTimeRaceWeekendEntryListSort,
	},
	{
		Name:     "Fewest Collisions",
		Key:      "fewest_collisions",
		SortFunc: FewestCollisionsRaceWeekendEntryListSort,
	},
	{
		Name:     "Fewest Cuts",
		Key:      "fewest_cuts",
		SortFunc: FewestCutsRaceWeekendEntryListSort,
	},
	{
		Name:     "Safety (Collisions then Cuts)",
		Key:      "safety",
		SortFunc: SafetyRaceWeekendEntryListSort,
	},
	{
		Name:     "Random",
		Key:      "random",
		SortFunc: RandomRaceWeekendEntryListSort,
	},
	{
		Name:     "Alphabetical (Using Driver Name)",
		Key:      "alphabetical",
		SortFunc: AlphabeticalRaceWeekendEntryListSort,
>>>>>>> origin/multiserver2
	},
}

func GetRaceWeekendEntryListSort(key string) RaceWeekendEntryListSorter {
	for _, sorter := range RaceWeekendEntryListSorters {
		if sorter.Key == key {
<<<<<<< HEAD
			return PerClassSort(sorter.Sorter)
		}
	}

	return PerClassSort(RaceWeekendEntryListSortFunc(UnchangedRaceWeekendEntryListSort))
}

func PerClassSort(sorter RaceWeekendEntryListSorter) RaceWeekendEntryListSorter {
	return RaceWeekendEntryListSortFunc(func(rw *RaceWeekend, session *RaceWeekendSession, allEntrants []*RaceWeekendSessionEntrant, filter *RaceWeekendSessionToSessionFilter) error {
		if _, isChampionshipClassSort := sorter.(*ChampionshipClassSort); isChampionshipClassSort && rw.HasLinkedChampionship() && rw.Championship != nil {
			// per championship class sort is a stable non-results based sort. If that has been selected, don't run this sorting function at all.
			return sorter.Sort(rw, session, allEntrants, filter)
		}

		classMap := make(map[uuid.UUID]bool)
=======
			return PerClassSort(sorter.SortFunc)
		}
	}

	return PerClassSort(UnchangedRaceWeekendEntryListSort)
}

func PerClassSort(sorter RaceWeekendEntryListSorter) RaceWeekendEntryListSorter {
	return func(session *RaceWeekendSession, allEntrants []*RaceWeekendSessionEntrant) error {
>>>>>>> origin/multiserver2
		fastestLapForClass := make(map[uuid.UUID]int)
		entrantsForClass := make(map[uuid.UUID][]*RaceWeekendSessionEntrant)

		for _, entrant := range allEntrants {
<<<<<<< HEAD
			classMap[entrant.EntrantResult.ClassID] = true

=======
>>>>>>> origin/multiserver2
			fastestLap, ok := fastestLapForClass[entrant.EntrantResult.ClassID]

			if entrant.EntrantResult.BestLap > 0 {
				if !ok || (ok && entrant.EntrantResult.BestLap < fastestLap) {
					fastestLapForClass[entrant.EntrantResult.ClassID] = entrant.EntrantResult.BestLap
				}
			}

			entrantsForClass[entrant.EntrantResult.ClassID] = append(entrantsForClass[entrant.EntrantResult.ClassID], entrant)
		}

		var classes []uuid.UUID

<<<<<<< HEAD
		for class := range classMap {
			classes = append(classes, class)
		}

		if len(fastestLapForClass) == len(classes) {
			// sort each class by the fastest lap in that class
			sort.Slice(classes, func(i, j int) bool {
				return fastestLapForClass[classes[i]] < fastestLapForClass[classes[j]]
			})
		} else if session.IsBase() {
			// base sessions will have no lap data. just sort them by class ID,
			// the same way every time, so that entrant splits can be consistent
			sort.Slice(classes, func(i, j int) bool {
				return classes[i].String() < classes[j].String()
			})
		}
=======
		for class := range fastestLapForClass {
			classes = append(classes, class)
		}

		// sort each class by the fastest lap in that class
		sort.Slice(classes, func(i, j int) bool {
			return fastestLapForClass[classes[i]] < fastestLapForClass[classes[j]]
		})
>>>>>>> origin/multiserver2

		lastStartPos := 0

		for _, class := range classes {
			entrants := entrantsForClass[class]

<<<<<<< HEAD
			err := sorter.Sort(rw, session, entrants, filter)
=======
			err := sorter(session, entrants)
>>>>>>> origin/multiserver2

			if err != nil {
				return err
			}

			reverseEntrants(session.NumEntrantsToReverse, entrants)

<<<<<<< HEAD
			if _, isChampionshipOrderSort := sorter.(*ChampionshipStandingsOrderEntryListSort); isChampionshipOrderSort && rw.HasLinkedChampionship() && rw.Championship != nil {
				sortDriversWithNoChampionshipRacesToBackOfGrid(rw.Championship, entrants)
			} else {
				sortDriversWithNoTimeToBackOfGrid(entrants)
			}
=======
			sortDriversWithNoTimeToBackOfGrid(entrants)
>>>>>>> origin/multiserver2

			for index, entrant := range entrants {
				allEntrants[lastStartPos+index] = entrant
			}

			lastStartPos += len(entrants)
		}

		return nil
<<<<<<< HEAD
	})
=======
	}
>>>>>>> origin/multiserver2
}

func sortDriversWithNoTimeToBackOfGrid(entrants []*RaceWeekendSessionEntrant) {
	sort.SliceStable(entrants, func(i, j int) bool {
		entrantI, entrantJ := entrants[i], entrants[j]

		if entrantI.EntrantResult.TotalTime == 0 {
			return false
		}

		if entrantJ.EntrantResult.TotalTime == 0 {
			return true
		}

		return i < j
	})
}

<<<<<<< HEAD
func UnchangedRaceWeekendEntryListSort(_ *RaceWeekend, _ *RaceWeekendSession, _ []*RaceWeekendSessionEntrant, _ *RaceWeekendSessionToSessionFilter) error {
	return nil // do nothing
}

func FastestLapRaceWeekendEntryListSort(rw *RaceWeekend, session *RaceWeekendSession, entrants []*RaceWeekendSessionEntrant, _ *RaceWeekendSessionToSessionFilter) error {
	sort.Slice(entrants, func(i, j int) bool {
		entrantI, entrantJ := entrants[i], entrants[j]

		return lessBestLapTime(rw, session, entrantI, entrantJ)
=======
func UnchangedRaceWeekendEntryListSort(_ *RaceWeekendSession, _ []*RaceWeekendSessionEntrant) error {
	return nil // do nothing
}

func FastestLapRaceWeekendEntryListSort(session *RaceWeekendSession, entrants []*RaceWeekendSessionEntrant) error {
	sort.Slice(entrants, func(i, j int) bool {
		entrantI, entrantJ := entrants[i], entrants[j]

		return lessBestLapTime(session, entrantI, entrantJ)
>>>>>>> origin/multiserver2
	})

	return nil
}

<<<<<<< HEAD
func FastestResultsFileRaceWeekendEntryListSort(_ *RaceWeekend, _ *RaceWeekendSession, entrants []*RaceWeekendSessionEntrant, filter *RaceWeekendSessionToSessionFilter) error {

	if filter == nil {
		return nil
	}

	bestDriverLaps := make(map[string]int)

	for _, resultFile := range filter.AvailableResultsForSorting {
		result, err := LoadResult(resultFile + ".json")

		if err != nil {
			return err
		}

		for _, driverResult := range result.Result {
			if driverResult.BestLap != 0 && (driverResult.BestLap < bestDriverLaps[driverResult.DriverGUID] || bestDriverLaps[driverResult.DriverGUID] == 0) {
				bestDriverLaps[driverResult.DriverGUID] = driverResult.BestLap
			}
		}
	}

	sort.Slice(entrants, func(i, j int) bool {
		entrantI, entrantJ := entrants[i], entrants[j]

		return lessBestLapTimeInResults(bestDriverLaps, entrantI, entrantJ)
	})

	return nil
}

func NumberResultsFileRaceWeekendEntryListSort(_ *RaceWeekend, _ *RaceWeekendSession, entrants []*RaceWeekendSessionEntrant, filter *RaceWeekendSessionToSessionFilter) error {

	if filter == nil {
		return nil
	}

	numDriverLaps := make(map[string]int)

	for _, resultFile := range filter.AvailableResultsForSorting {
		result, err := LoadResult(resultFile + ".json")

		if err != nil {
			return err
		}

		for _, sessionResult := range result.Result {
			numDriverLaps[sessionResult.DriverGUID] += result.GetNumLaps(sessionResult.DriverGUID, sessionResult.CarModel)
		}
	}

	sort.Slice(entrants, func(i, j int) bool {
		return numDriverLaps[entrants[i].Car.Driver.GUID] > numDriverLaps[entrants[j].Car.Driver.GUID]
	})

	return nil
}

func TotalRaceTimeRaceWeekendEntryListSort(rw *RaceWeekend, session *RaceWeekendSession, entrants []*RaceWeekendSessionEntrant, _ *RaceWeekendSessionToSessionFilter) error {
	sort.Slice(entrants, func(i, j int) bool {
		entrantI, entrantJ := entrants[i], entrants[j]

		return lessTotalEntrantTime(rw, session, entrantI, entrantJ)
=======
func TotalRaceTimeRaceWeekendEntryListSort(session *RaceWeekendSession, entrants []*RaceWeekendSessionEntrant) error {
	sort.Slice(entrants, func(i, j int) bool {
		entrantI, entrantJ := entrants[i], entrants[j]

		return lessTotalEntrantTime(session, entrantI, entrantJ)
>>>>>>> origin/multiserver2
	})

	return nil
}

<<<<<<< HEAD
func lessTotalEntrantTime(_ *RaceWeekend, _ *RaceWeekendSession, entrantI, entrantJ *RaceWeekendSessionEntrant) bool {
	if entrantI.SessionResults.GetNumLaps(entrantI.Car.Driver.GUID, entrantI.Car.Model) == entrantJ.SessionResults.GetNumLaps(entrantJ.Car.Driver.GUID, entrantJ.Car.Model) {
		// drivers have completed the same number of laps, so compare their total time
		entrantITime := entrantI.SessionResults.GetTime(entrantI.EntrantResult.TotalTime, entrantI.Car.Driver.GUID, entrantI.Car.Model, true)
		entrantJTime := entrantJ.SessionResults.GetTime(entrantJ.EntrantResult.TotalTime, entrantJ.Car.Driver.GUID, entrantJ.Car.Model, true)

		return entrantITime < entrantJTime
	}

	return entrantI.SessionResults.GetNumLaps(entrantI.Car.Driver.GUID, entrantI.Car.Model) > entrantJ.SessionResults.GetNumLaps(entrantJ.Car.Driver.GUID, entrantJ.Car.Model)
}

func lessBestLapTime(_ *RaceWeekend, _ *RaceWeekendSession, entrantI, entrantJ *RaceWeekendSessionEntrant) bool {
=======
func lessTotalEntrantTime(session *RaceWeekendSession, entrantI, entrantJ *RaceWeekendSessionEntrant) bool {
	if entrantI.SessionResults.GetNumLaps(entrantI.Car.Driver.GUID) == entrantJ.SessionResults.GetNumLaps(entrantJ.Car.Driver.GUID) {
		// drivers have completed the same number of laps, so compare their total time
		entrantITime := entrantI.SessionResults.GetTime(entrantI.EntrantResult.TotalTime, entrantI.Car.Driver.GUID, true)
		entrantJTime := entrantJ.SessionResults.GetTime(entrantJ.EntrantResult.TotalTime, entrantJ.Car.Driver.GUID, true)

		return entrantITime < entrantJTime
	} else {
		return entrantI.SessionResults.GetNumLaps(entrantI.Car.Driver.GUID) > entrantJ.SessionResults.GetNumLaps(entrantJ.Car.Driver.GUID)
	}
}

func lessBestLapTime(session *RaceWeekendSession, entrantI, entrantJ *RaceWeekendSessionEntrant) bool {
>>>>>>> origin/multiserver2
	if entrantI.EntrantResult.BestLap == 0 {
		// entrantI has a zero best lap. they must be not-less than J
		return false
	}

	if entrantJ.EntrantResult.BestLap == 0 {
		// entrantJ has a zero best lap. entrantI must be less than J
		return true
	}

	if entrantI.EntrantResult.BestLap == entrantJ.EntrantResult.BestLap {
		// if equal, compare safety
<<<<<<< HEAD
		entrantICrashes := entrantI.SessionResults.GetCrashes(entrantI.Car.Driver.GUID, entrantI.Car.Model)
		entrantJCrashes := entrantJ.SessionResults.GetCrashes(entrantJ.Car.Driver.GUID, entrantJ.Car.Model)

		if entrantICrashes == entrantJCrashes {
			return entrantI.SessionResults.GetCuts(entrantI.Car.Driver.GUID, entrantI.Car.Model) < entrantJ.SessionResults.GetCuts(entrantJ.Car.Driver.GUID, entrantJ.Car.Model)
=======
		entrantICrashes := entrantI.SessionResults.GetCrashes(entrantI.Car.Driver.GUID)
		entrantJCrashes := entrantJ.SessionResults.GetCrashes(entrantJ.Car.Driver.GUID)

		if entrantICrashes == entrantJCrashes {
			return entrantI.SessionResults.GetCuts(entrantI.Car.Driver.GUID) < entrantJ.SessionResults.GetCuts(entrantJ.Car.Driver.GUID)
>>>>>>> origin/multiserver2
		}

		return entrantICrashes < entrantJCrashes
	}

	return entrantI.EntrantResult.BestLap < entrantJ.EntrantResult.BestLap
}

<<<<<<< HEAD
func lessBestLapTimeInResults(bestDriverLaps map[string]int, entrantI, entrantJ *RaceWeekendSessionEntrant) bool {
	if bestDriverLaps[entrantI.Car.Driver.GUID] == 0 {
		return false
	}

	if bestDriverLaps[entrantJ.Car.Driver.GUID] == 0 {
		return true
	}

	return bestDriverLaps[entrantI.Car.Driver.GUID] < bestDriverLaps[entrantJ.Car.Driver.GUID]
}

func FewestCollisionsRaceWeekendEntryListSort(rw *RaceWeekend, session *RaceWeekendSession, entrants []*RaceWeekendSessionEntrant, _ *RaceWeekendSessionToSessionFilter) error {
	sort.Slice(entrants, func(i, j int) bool {
		entrantI, entrantJ := entrants[i], entrants[j]
		entrantICrashes := entrantI.SessionResults.GetCrashes(entrantI.Car.Driver.GUID, entrantI.Car.Model)
		entrantJCrashes := entrantJ.SessionResults.GetCrashes(entrantJ.Car.Driver.GUID, entrantJ.Car.Model)

		if entrantICrashes == entrantJCrashes {
			if session.SessionType() == SessionTypeRace {
				return lessTotalEntrantTime(rw, session, entrantI, entrantJ)
			}

			return lessBestLapTime(rw, session, entrantI, entrantJ)
=======
func FewestCollisionsRaceWeekendEntryListSort(session *RaceWeekendSession, entrants []*RaceWeekendSessionEntrant) error {
	sort.Slice(entrants, func(i, j int) bool {
		entrantI, entrantJ := entrants[i], entrants[j]
		entrantICrashes := entrantI.SessionResults.GetCrashes(entrantI.Car.Driver.GUID)
		entrantJCrashes := entrantJ.SessionResults.GetCrashes(entrantJ.Car.Driver.GUID)

		if entrantICrashes == entrantJCrashes {
			if session.SessionType() == SessionTypeRace {
				return lessTotalEntrantTime(session, entrantI, entrantJ)
			} else {
				return lessBestLapTime(session, entrantI, entrantJ)
			}
>>>>>>> origin/multiserver2
		}

		return entrantICrashes < entrantJCrashes
	})

	return nil
}

<<<<<<< HEAD
func FewestCutsRaceWeekendEntryListSort(rw *RaceWeekend, session *RaceWeekendSession, entrants []*RaceWeekendSessionEntrant, _ *RaceWeekendSessionToSessionFilter) error {
	sort.Slice(entrants, func(i, j int) bool {
		entrantI, entrantJ := entrants[i], entrants[j]
		entrantICuts := entrantI.SessionResults.GetCuts(entrantI.Car.Driver.GUID, entrantI.Car.Model)
		entrantJCuts := entrantJ.SessionResults.GetCuts(entrantJ.Car.Driver.GUID, entrantJ.Car.Model)

		if entrantICuts == entrantJCuts {
			if session.SessionType() == SessionTypeRace {
				return lessTotalEntrantTime(rw, session, entrantI, entrantJ)
			}

			return lessBestLapTime(rw, session, entrantI, entrantJ)
=======
func FewestCutsRaceWeekendEntryListSort(session *RaceWeekendSession, entrants []*RaceWeekendSessionEntrant) error {
	sort.Slice(entrants, func(i, j int) bool {
		entrantI, entrantJ := entrants[i], entrants[j]
		entrantICuts := entrantI.SessionResults.GetCuts(entrantI.Car.Driver.GUID)
		entrantJCuts := entrantJ.SessionResults.GetCuts(entrantJ.Car.Driver.GUID)

		if entrantICuts == entrantJCuts {
			if session.SessionType() == SessionTypeRace {
				return lessTotalEntrantTime(session, entrantI, entrantJ)
			} else {
				return lessBestLapTime(session, entrantI, entrantJ)
			}
>>>>>>> origin/multiserver2
		}

		return entrantICuts < entrantJCuts
	})

	return nil
}

<<<<<<< HEAD
func SafetyRaceWeekendEntryListSort(rw *RaceWeekend, session *RaceWeekendSession, entrants []*RaceWeekendSessionEntrant, _ *RaceWeekendSessionToSessionFilter) error {
	sort.Slice(entrants, func(i, j int) bool {
		entrantI, entrantJ := entrants[i], entrants[j]
		entrantICrashes := entrantI.SessionResults.GetCrashes(entrantI.Car.Driver.GUID, entrantI.Car.Model)
		entrantJCrashes := entrantJ.SessionResults.GetCrashes(entrantJ.Car.Driver.GUID, entrantJ.Car.Model)
		entrantICuts := entrantI.SessionResults.GetCuts(entrantI.Car.Driver.GUID, entrantI.Car.Model)
		entrantJCuts := entrantJ.SessionResults.GetCuts(entrantJ.Car.Driver.GUID, entrantJ.Car.Model)
=======
func SafetyRaceWeekendEntryListSort(session *RaceWeekendSession, entrants []*RaceWeekendSessionEntrant) error {
	sort.Slice(entrants, func(i, j int) bool {
		entrantI, entrantJ := entrants[i], entrants[j]
		entrantICrashes := entrantI.SessionResults.GetCrashes(entrantI.Car.Driver.GUID)
		entrantJCrashes := entrantJ.SessionResults.GetCrashes(entrantJ.Car.Driver.GUID)
		entrantICuts := entrantI.SessionResults.GetCuts(entrantI.Car.Driver.GUID)
		entrantJCuts := entrantJ.SessionResults.GetCuts(entrantJ.Car.Driver.GUID)
>>>>>>> origin/multiserver2

		if entrantICrashes == entrantJCrashes {
			if entrantICuts == entrantJCuts {
				if session.SessionType() == SessionTypeRace {
<<<<<<< HEAD
					return lessTotalEntrantTime(rw, session, entrantI, entrantJ)
				}

				return lessBestLapTime(rw, session, entrantI, entrantJ)
			}

			return entrantICuts < entrantJCuts
=======
					return lessTotalEntrantTime(session, entrantI, entrantJ)
				} else {
					return lessBestLapTime(session, entrantI, entrantJ)
				}
			} else {
				return entrantICuts < entrantJCuts
			}
>>>>>>> origin/multiserver2
		}

		return entrantICrashes < entrantJCrashes
	})

	return nil
}

<<<<<<< HEAD
func RandomRaceWeekendEntryListSort(_ *RaceWeekend, _ *RaceWeekendSession, entrants []*RaceWeekendSessionEntrant, _ *RaceWeekendSessionToSessionFilter) error {
=======
func RandomRaceWeekendEntryListSort(_ *RaceWeekendSession, entrants []*RaceWeekendSessionEntrant) error {
>>>>>>> origin/multiserver2
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	r.Shuffle(len(entrants), func(i, j int) {
		entrants[i], entrants[j] = entrants[j], entrants[i]
	})

	return nil
}

<<<<<<< HEAD
func AlphabeticalRaceWeekendEntryListSort(_ *RaceWeekend, _ *RaceWeekendSession, entrants []*RaceWeekendSessionEntrant, _ *RaceWeekendSessionToSessionFilter) error {
=======
func AlphabeticalRaceWeekendEntryListSort(_ *RaceWeekendSession, entrants []*RaceWeekendSessionEntrant) error {
>>>>>>> origin/multiserver2
	sort.Slice(entrants, func(i, j int) bool {
		entrantI, entrantJ := entrants[i], entrants[j]

		return entrantI.Car.Driver.Name < entrantJ.Car.Driver.Name
	})

	return nil
}
<<<<<<< HEAD

type ChampionshipStandingsOrderEntryListSort struct{}

func (ChampionshipStandingsOrderEntryListSort) Sort(rw *RaceWeekend, _ *RaceWeekendSession, entrants []*RaceWeekendSessionEntrant, _ *RaceWeekendSessionToSessionFilter) error {
	if !rw.HasLinkedChampionship() || rw.Championship == nil {
		return nil
	}

	if len(entrants) == 0 {
		return nil
	}

	class := entrants[0].ChampionshipClass(rw)
	standings := class.Standings(rw.Championship, rw.Championship.Events)

	sort.Slice(entrants, func(i, j int) bool {
		entrantI, entrantJ := entrants[i], entrants[j]

		iPos := len(standings)
		jPos := len(standings)

		for i, standing := range standings {
			if standing.Car.GetGUID() == entrantI.GetEntrant().GUID {
				iPos = i
			}

			if standing.Car.GetGUID() == entrantJ.GetEntrant().GUID {
				jPos = i
			}
		}

		return iPos < jPos
	})

	return nil
}

type ChampionshipClassSort struct{}

func (ChampionshipClassSort) Sort(rw *RaceWeekend, _ *RaceWeekendSession, entrants []*RaceWeekendSessionEntrant, _ *RaceWeekendSessionToSessionFilter) error {
	if !rw.HasLinkedChampionship() || rw.Championship == nil {
		return nil
	}

	if len(entrants) == 0 {
		return nil
	}

	sort.Slice(entrants, func(i, j int) bool {
		return entrants[i].ChampionshipClass(rw).ID.String() < entrants[j].ChampionshipClass(rw).ID.String()
	})

	return nil
}

func sortDriversWithNoChampionshipRacesToBackOfGrid(championship *Championship, entrants []*RaceWeekendSessionEntrant) {
	sort.SliceStable(entrants, func(i, j int) bool {
		entrantI, entrantJ := entrants[i], entrants[j]

		entrantIAttendance := championship.EntrantAttendance(entrantI.GetEntrant().GUID)
		entrantJAttendance := championship.EntrantAttendance(entrantJ.GetEntrant().GUID)

		if entrantIAttendance == 0 {
			return false
		}

		if entrantJAttendance == 0 {
			return true
		}

		return i < j
	})
}
=======
>>>>>>> origin/multiserver2
