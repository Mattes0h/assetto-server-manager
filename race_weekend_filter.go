package servermanager

<<<<<<< HEAD
import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/cj123/ini"
	"github.com/sirupsen/logrus"
)

=======
>>>>>>> origin/multiserver2
type FilterError string

func (f FilterError) Error() string {
	return string(f)
}

type RaceWeekendSessionToSessionFilter struct {
<<<<<<< HEAD
	// IsPreview indicates that the Filter is for preview only, and will not actually affect a starting grid.
	IsPreview bool

	// ResultStart is the beginning of the split from the previous session's result
	ResultStart int
	// ResultEnd is the end of the split from the previous session's result
	ResultEnd int

	// NumEntrantsToReverse defines how many entrants to reverse. -1 indicates all, 0 indicates none, or N entrants.
	NumEntrantsToReverse int

	// EntryListStart is where to place the entrants in the starting grid of the next session
	EntryListStart int

	// SortType defines how the entrants are sorted
	SortType string

	// ForceUseTyreFromFastestLap forces drivers to start on the same tyre compound as the tyre compound that
	// they achieved their fastest lap on in the previous session
	ForceUseTyreFromFastestLap bool

	// AvailableResultsForSorting are the results files to be used for sorting the entrants.
	AvailableResultsForSorting []string

	// ManualDriverSelection indicates that drivers are picked manually from the above results file.
	ManualDriverSelection bool

	// SelectedDriverGUIDs is a list of the currently selected driver GUIDs.
	SelectedDriverGUIDs []string
=======
	ResultStart int
	ResultEnd   int

	NumEntrantsToReverse int

	EntryListStart int

	SortType string
>>>>>>> origin/multiserver2
}

func reverseEntrants(numToReverse int, entrants []*RaceWeekendSessionEntrant) {
	if numToReverse == 0 {
		return
	}

	if numToReverse > len(entrants) {
		numToReverse = len(entrants)
	}

	var toReverse []*RaceWeekendSessionEntrant

	if numToReverse > 0 {
		toReverse = entrants[:numToReverse]
	} else {
		toReverse = entrants
	}

	for i := len(toReverse)/2 - 1; i >= 0; i-- {
		opp := len(toReverse) - 1 - i
		toReverse[i], toReverse[opp] = toReverse[opp], toReverse[i]
	}

	for i := 0; i < len(toReverse); i++ {
		entrants[i] = toReverse[i]
	}
}

// Filter takes a set of RaceWeekendSessionEntrants formed by the results of the parent session and filters them into a child session entry list.
<<<<<<< HEAD
func (f RaceWeekendSessionToSessionFilter) Filter(raceWeekend *RaceWeekend, parentSession, childSession *RaceWeekendSession, parentSessionResults []*RaceWeekendSessionEntrant, childSessionEntryList *RaceWeekendEntryList) error {
	if parentSession.Completed() || childSession.IsBase() {
=======
func (f RaceWeekendSessionToSessionFilter) Filter(parentSession, childSession *RaceWeekendSession, parentSessionResults []*RaceWeekendSessionEntrant, childSessionEntryList *RaceWeekendEntryList) error {
	if parentSession.Completed() {
>>>>>>> origin/multiserver2
		sorter := GetRaceWeekendEntryListSort(f.SortType)

		parentSession.NumEntrantsToReverse = f.NumEntrantsToReverse

		// race weekend session is completed and has a valid sorter, use it to sort results before filtering.
<<<<<<< HEAD
		if err := sorter.Sort(raceWeekend, parentSession, parentSessionResults, &f); err != nil {
=======
		if err := sorter(parentSession, parentSessionResults); err != nil {
>>>>>>> origin/multiserver2
			return err
		}
	}

<<<<<<< HEAD
	entryListStart := f.EntryListStart - 1

	var split []*RaceWeekendSessionEntrant

	if f.ManualDriverSelection {
		for _, driverGUID := range f.SelectedDriverGUIDs {
			for _, entrant := range parentSessionResults {
				if entrant.Car.GetGUID() == driverGUID {
					split = append(split, entrant)
					break
				}
			}
		}
	} else {
		resultStart, resultEnd := f.ResultStart, f.ResultEnd

		resultStart--

		if resultStart > len(parentSessionResults) {
			return nil
		}

		if resultEnd > len(parentSessionResults) {
			resultEnd = len(parentSessionResults)
		}

		split = parentSessionResults[resultStart:resultEnd]
	}

=======
	resultStart, resultEnd, entryListStart := f.ResultStart, f.ResultEnd, f.EntryListStart

	resultStart--
	entryListStart--

	if resultStart > len(parentSessionResults) {
		return nil
	}

	if resultEnd > len(parentSessionResults) {
		resultEnd = len(parentSessionResults)
	}

	split := parentSessionResults[resultStart:resultEnd]

>>>>>>> origin/multiserver2
	if !parentSession.Completed() {
		reverseEntrants(f.NumEntrantsToReverse, split)
	}

	splitIndex := 0

<<<<<<< HEAD
	for pitBox := entryListStart; pitBox < entryListStart+len(split); pitBox++ {
		entrant := split[splitIndex]
		entrant.SessionID = parentSession.ID

		if !f.IsPreview && parentSession.Completed() && f.ForceUseTyreFromFastestLap {
			// find the tyre from the entrants fastest lap
			fastestLap := entrant.SessionResults.GetDriversFastestLap(entrant.Car.GetGUID(), entrant.Car.GetCar())

			if fastestLap == nil {
				logrus.Warnf("could not find fastest lap for entrant %s (%s). will not lock their tyre choice.", entrant.Car.GetName(), entrant.Car.GetGUID())
			} else {
				err := raceWeekend.buildLockedTyreSetup(childSession, entrant, fastestLap)

				if err != nil {
					logrus.WithError(err).Errorf("could not build locked tyre setup for entrant %s (%s)", entrant.Car.GetName(), entrant.Car.GetGUID())
				}
			}
		}

=======
	for pitBox := entryListStart; pitBox < entryListStart+(resultEnd-resultStart); pitBox++ {
		entrant := split[splitIndex]
		entrant.SessionID = parentSession.ID

>>>>>>> origin/multiserver2
		childSessionEntryList.AddInPitBox(entrant, pitBox)

		splitIndex++
	}

	return nil
}
<<<<<<< HEAD

const lockedTyreSetupFolder = "server_manager_locked_tyres"

func (rw *RaceWeekend) buildLockedTyreSetup(session *RaceWeekendSession, entrant *RaceWeekendSessionEntrant, fastestLap *SessionLap) error {
	tyreIndex, err := findTyreIndex(entrant.Car.Model, fastestLap.Tyre, session.RaceConfig)

	if err != nil {
		return err
	}

	entryList := rw.GetEntryList()

	var setup *ini.File

	for _, raceWeekendEntrant := range entryList {
		if raceWeekendEntrant.GUID == entrant.Car.GetGUID() && raceWeekendEntrant.FixedSetup != "" {
			setup, err = ini.Load(filepath.Join(ServerInstallPath, "setups", raceWeekendEntrant.FixedSetup))

			if err != nil {
				return err
			}

			break
		}
	}

	if setup == nil {
		// no fixed setup was specified
		// write out a temp ini setup file for this car + player.
		setup = ini.NewFile([]ini.DataSource{nil}, ini.LoadOptions{
			IgnoreInlineComment: true,
		})

		_, err = setup.NewSection("DEFAULT")

		if err != nil {
			return err
		}

		car, err := setup.NewSection("CAR")

		if err != nil {
			return err
		}

		_, err = car.NewKey("MODEL", entrant.Car.Model)

		if err != nil {
			return err
		}
	}

	tyres, err := setup.NewSection("TYRES")

	if err != nil {
		return err
	}

	_, err = tyres.NewKey("VALUE", fmt.Sprintf("%d", tyreIndex))

	if err != nil {
		return err
	}

	raceWeekendSection, err := setup.NewSection("RACE_WEEKEND")

	if err != nil {
		return err
	}

	_, err = raceWeekendSection.NewKey("ID", rw.ID.String())

	if err != nil {
		return err
	}

	_, err = raceWeekendSection.NewKey("SESSION_ID", entrant.SessionID.String())

	if err != nil {
		return err
	}

	setupFilePath := filepath.Join(entrant.Car.Model, lockedTyreSetupFolder, fmt.Sprintf("race_weekend_session_%s_%s.ini", entrant.Car.GetGUID(), entrant.SessionID.String()))

	fullSaveFilepath := filepath.Join(ServerInstallPath, "setups", setupFilePath)

	if err := os.MkdirAll(filepath.Dir(fullSaveFilepath), 0755); err != nil {
		return err
	}

	if err := setup.SaveTo(fullSaveFilepath); err != nil {
		return err
	}

	entrant.OverrideSetupFile = setupFilePath

	return nil
}
=======
>>>>>>> origin/multiserver2
