import {
    RaceControl as RaceControlData,
    RaceControlDriverMapRaceControlDriver as Driver,
    RaceControlDriverMapRaceControlDriverSessionCarInfo as SessionCarInfo
} from "./models/RaceControl";

import {CarUpdate, CarUpdateVec} from "./models/UDP";
import {randomColor} from "randomcolor/randomColor";
import {msToTime, prettifyName} from "./utils";
import moment from "moment";
import ClickEvent = JQuery.ClickEvent;

interface WSMessage {
    Message: any;
    EventType: number;
}

const EventCollisionWithCar = 10,
    EventCollisionWithEnv = 11,
    EventNewSession = 50,
    EventNewConnection = 51,
    EventConnectionClosed = 52,
    EventCarUpdate = 53,
    EventCarInfo = 54,
    EventEndSession = 55,
    EventVersion = 56,
    EventChat = 57,
    EventClientLoaded = 58,
    EventSessionInfo = 59,
    EventError = 60,
    EventLapCompleted = 73,
    EventClientEvent = 130,
    EventRaceControl = 200
;

interface SimpleCollision {
    WorldPos: CarUpdateVec
}

interface WebsocketHandler {
    handleWebsocketMessage(message: WSMessage): void;
}

export class RaceControl {
    private readonly liveMap: LiveMap = new LiveMap(this);
    private readonly liveTimings: LiveTimings = new LiveTimings(this, this.liveMap);
    private readonly $eventTitle: JQuery<HTMLHeadElement>;
    public status: RaceControlData;

    constructor() {
        this.$eventTitle = $("#event-title");
        this.status = new RaceControlData();

        if (!this.$eventTitle.length) {
            return;
        }

        // enable wide-mode
        $(".container").attr("class", "container-fluid");

        let ws = new WebSocket(((window.location.protocol === "https:") ? "wss://" : "ws://") + window.location.host + "/api/race-control");
        ws.onmessage = this.handleWebsocketMessage.bind(this);

        setInterval(this.showEventCompletion.bind(this), 1000);
    }

    private handleWebsocketMessage(ev: MessageEvent): void {
        let message = JSON.parse(ev.data) as WSMessage;

        if (!message) {
            return;
        }

        switch (message.EventType) {
            case EventRaceControl:
                this.status = new RaceControlData(message.Message);
                this.$eventTitle.text(RaceControl.getSessionType(this.status.SessionInfo.Type) + " at " + this.status.TrackInfo!.name);

                this.buildSessionInfo();
                break;
        }

        this.liveMap.handleWebsocketMessage(message);
        this.liveTimings.handleWebsocketMessage(message);
    }

    private static getSessionType(sessionIndex: number): string {
        switch (sessionIndex) {
            case 0:
                return "Booking";
            case 1:
                return "Practice";
            case 2:
                return "Qualifying";
            case 3:
                return "Race";
            default:
                return "Unknown session";
        }
    }

    private showEventCompletion() {
        let timeRemaining = "";

        // Get lap/laps or time/totalTime
        if (this.status.SessionInfo.Time > 0) {
            timeRemaining = msToTime(this.status.SessionInfo.Time * 60 * 1000 - moment.duration(moment().diff(this.status.SessionStartTime)).asMilliseconds(), false, false);
        } else if (this.status.SessionInfo.Laps > 0) {
            let lapsCompleted = 0;

            if (this.status.ConnectedDrivers && this.status.ConnectedDrivers.GUIDsInPositionalOrder.length > 0) {
                let driver = this.status.ConnectedDrivers.Drivers[this.status.ConnectedDrivers.GUIDsInPositionalOrder[0]];

                if (driver.TotalNumLaps > 0) {
                    lapsCompleted = driver.TotalNumLaps;
                }
            }

            timeRemaining = this.status.SessionInfo.Laps - lapsCompleted + " laps remaining";
        }

        let $raceTime = $("#race-time");
        $raceTime.text(timeRemaining);
    }

    private buildSessionInfo() {
        let $roadTempWrapper = $("#road-temp-wrapper");
        $roadTempWrapper.attr("style", "background-color: " + getColorForPercentage(this.status.SessionInfo.RoadTemp / 40));
        $roadTempWrapper.attr("data-original-title", "Road Temp: " + this.status.SessionInfo.RoadTemp + "°C");

        let $roadTempText = $("#road-temp-text");
        $roadTempText.text(this.status.SessionInfo.RoadTemp + "°C");

        let $ambientTempWrapper = $("#ambient-temp-wrapper");
        $ambientTempWrapper.attr("style", "background-color: " + getColorForPercentage(this.status.SessionInfo.AmbientTemp / 40));
        $ambientTempWrapper.attr("data-original-title", "Ambient Temp: " + this.status.SessionInfo.AmbientTemp + "°C");

        let $ambientTempText = $("#ambient-temp-text");
        $ambientTempText.text(this.status.SessionInfo.AmbientTemp + "°C");

        // @TODO only needs changing every new session or on init (i.e. when no this.status)
        let $currentWeather = $("#weatherImage");

        // Fix for sol weathers with time info in this format:
        // sol_05_Broken%20Clouds_type=18_time=0_mult=20_start=1551792960/preview.jpg
        let pathCorrected = this.status.SessionInfo.WeatherGraphics.split("_");

        for (let i = 0; i < pathCorrected.length; i++) {
            if (pathCorrected[i].indexOf("type=") !== -1) {
                pathCorrected.splice(i);
                break;
            }
        }

        let pathFinal = pathCorrected.join("_");

        $.get("/content/weather/" + pathFinal + "/preview.jpg").done(function () {
            // preview for skin exists
            $currentWeather.attr("src", "/content/weather/" + pathFinal + "/preview.jpg");
        }).fail(function () {
            // preview doesn't exist, load default fall back image
            $currentWeather.attr("src", "/static/img/no-preview-general.png");
        });

        $currentWeather.attr("alt", "Current Weather: " + prettifyName(this.status.SessionInfo.WeatherGraphics, false));
        $("#trackImage").attr("src", this.getTrackImageURL());

        $("#event-name").text(this.status.SessionInfo.Name);
        $("#event-type").text(RaceControl.getSessionType(this.status.SessionInfo.Type));
    }

    private getTrackImageURL(): string {
        if (!this.status) {
            return "";
        }

        const sessionInfo = this.status.SessionInfo;

        return "/content/tracks/" + sessionInfo.Track + "/ui" + (!!sessionInfo.TrackConfig ? "/" + sessionInfo.TrackConfig : "") + "/preview.png";
    }
}

class LiveMap implements WebsocketHandler {
    private mapImageHasLoaded: boolean = false;

    private readonly $map: JQuery<HTMLDivElement>;
    private readonly $trackMapImage: JQuery<HTMLImageElement> | undefined;
    private readonly raceControl: RaceControl;

    constructor(raceControl: RaceControl) {
        this.$map = $("#map");
        this.raceControl = raceControl;
        this.$trackMapImage = this.$map.find("img") as JQuery<HTMLImageElement>;

        $(window).on("resize", this.correctMapDimensions.bind(this));
    }

    // positional coordinate modifiers.
    private mapScaleMultiplier: number = 1;
    private trackScale: number = 1;
    private trackMargin: number = 0;
    private trackXOffset: number = 0;
    private trackZOffset: number = 0;

    // live map track dots
    private dots: Map<string, JQuery<HTMLElement>> = new Map<string, JQuery<HTMLElement>>();
    private maxRPMs: Map<string, number> = new Map<string, number>();

    public handleWebsocketMessage(message: WSMessage): void {
        switch (message.EventType) {
            case EventRaceControl:
                this.trackXOffset = this.raceControl.status.TrackMapData!.offset_x;
                this.trackZOffset = this.raceControl.status.TrackMapData!.offset_y;
                this.trackScale = this.raceControl.status.TrackMapData!.scale_factor;
                this.loadTrackMapImage();

                for (const connectedGUID in this.raceControl.status.ConnectedDrivers!.Drivers) {
                    const driver = this.raceControl.status.ConnectedDrivers!.Drivers[connectedGUID];

                    if (!this.dots.has(driver.CarInfo.DriverGUID)) {
                        // in the event that a user just loaded the race control page, place the
                        // already loaded dots onto the map
                        let $driverDot = this.buildDriverDot(driver.CarInfo).show();
                        this.dots.set(driver.CarInfo.DriverGUID, $driverDot);
                    }
                }

                break;

            case EventNewConnection:
                const connectedDriver = new SessionCarInfo(message.Message);
                this.dots.set(connectedDriver.DriverGUID, this.buildDriverDot(connectedDriver));

                break;

            case EventClientLoaded:
                let carID = message.Message as number;

                if (!this.raceControl.status!.CarIDToGUID.hasOwnProperty(carID)) {
                    return;
                }

                // find the guid for this car ID:
                this.dots.get(this.raceControl.status!.CarIDToGUID[carID])!.show();

                break;

            case EventConnectionClosed:
                const disconnectedDriver = new SessionCarInfo(message.Message);
                const $dot = this.dots.get(disconnectedDriver.DriverGUID);

                if ($dot) {
                    $dot.hide();
                    this.dots.delete(disconnectedDriver.DriverGUID);
                }

                break;
            case EventCarUpdate:
                const update = new CarUpdate(message.Message);

                if (!this.raceControl.status!.CarIDToGUID.hasOwnProperty(update.CarID)) {
                    return;
                }

                // find the guid for this car ID:
                const driverGUID = this.raceControl.status!.CarIDToGUID[update.CarID];

                let $myDot = this.dots.get(driverGUID);
                let dotPos = this.translateToTrackCoordinate(update.Pos);

                $myDot!.css({
                    "left": dotPos.X,
                    "top": dotPos.Z,
                });

                let speed = Math.floor(Math.sqrt((Math.pow(update.Velocity.X, 2) + Math.pow(update.Velocity.Z, 2))) * 3.6);

                let maxRPM = this.maxRPMs.get(driverGUID);

                if (!maxRPM) {
                    maxRPM = 0;
                }

                if (update.EngineRPM > maxRPM) {
                    maxRPM = update.EngineRPM;
                    this.maxRPMs.set(driverGUID, update.EngineRPM);
                }

                let $rpmGaugeOuter = $("<div class='rpm-outer'></div>");
                let $rpmGaugeInner = $("<div class='rpm-inner'></div>");

                $rpmGaugeInner.css({
                    'width': ((update.EngineRPM / maxRPM) * 100).toFixed(0) + "%",
                    'background': randomColorForDriver(driverGUID),
                });

                $rpmGaugeOuter.append($rpmGaugeInner);
                $myDot!.find(".info").text(speed + "Km/h " + (update.Gear - 1));
                $myDot!.find(".info").append($rpmGaugeOuter);
                break;

            case EventNewSession:
                this.loadTrackMapImage();

                break;

            case EventCollisionWithCar:
            case EventCollisionWithEnv:
                let collisionData = message.Message as SimpleCollision;

                let collisionMapPoint = this.translateToTrackCoordinate(collisionData.WorldPos);

                let $collision = $("<div class='collision' />").css({
                    'left': collisionMapPoint.X,
                    'top': collisionMapPoint.Z,
                });

                $collision.appendTo(this.$map);

                break;
        }
    }

    private translateToTrackCoordinate(vec: CarUpdateVec): CarUpdateVec {
        const out = new CarUpdateVec();

        out.X = ((vec.X + this.trackXOffset + this.trackMargin) / this.trackScale) * this.mapScaleMultiplier;
        out.Z = ((vec.Z + this.trackZOffset + this.trackMargin) / this.trackScale) * this.mapScaleMultiplier;

        return out;
    }

    private buildDriverDot(driverData: SessionCarInfo): JQuery<HTMLElement> {
        if (this.dots.has(driverData.DriverGUID)) {
            return this.dots.get(driverData.DriverGUID)!;
        }

        const $driverName = $("<span class='name'/>").text(getAbbreviation(driverData.DriverName));
        const $info = $("<span class='info'/>").text("0").hide();

        const $dot = $("<div class='dot' style='background: " + randomColorForDriver(driverData.DriverGUID) + "'/>").append($driverName, $info).hide().appendTo(this.$map);

        this.dots.set(driverData.DriverGUID, $dot);

        return $dot;
    }

    private getTrackMapURL(): string {
        if (!this.raceControl.status) {
            return "";
        }

        const sessionInfo = this.raceControl.status.SessionInfo;

        return "/content/tracks/" + sessionInfo.Track + (!!sessionInfo.TrackConfig ? "/" + sessionInfo.TrackConfig : "") + "/map.png";
    }

    trackImage: HTMLImageElement = new Image();

    private loadTrackMapImage(): void {
        const trackURL = this.getTrackMapURL();

        this.trackImage.onload = () => {
            this.$trackMapImage!.attr({
                "src": trackURL,
            });

            this.mapImageHasLoaded = true;
            this.correctMapDimensions();
        };

        this.trackImage.src = trackURL
    }

    private static mapRotationRatio: number = 1.07;

    private correctMapDimensions(): void {
        if (!this.trackImage || !this.$trackMapImage || !this.mapImageHasLoaded) {
            return;
        }

        if (this.trackImage.height / this.trackImage.width > LiveMap.mapRotationRatio) {
            // rotate the map
            this.$map.addClass("rotated");

            this.$trackMapImage.css({
                'max-height': this.$trackMapImage.closest(".map-container").width()!,
                'max-width': 'auto'
            });

            this.mapScaleMultiplier = this.$trackMapImage.width()! / this.trackImage.width;

            this.$map.closest(".map-container").css({
                'max-height': (this.trackImage.width * this.mapScaleMultiplier) + 20,
            });

            this.$map.css({
                'max-width': (this.trackImage.width * this.mapScaleMultiplier) + 20,
            });
        } else {
            // un-rotate the map
            this.$map.removeClass("rotated").css({
                'max-height': 'inherit',
                'max-width': '100%',
            });

            this.$map.closest(".map-container").css({
                'max-height': 'auto',
            });

            this.$trackMapImage.css({
                'max-height': 'inherit',
                'max-width': '100%'
            });

            this.mapScaleMultiplier = this.$trackMapImage.width()! / this.trackImage.width;
        }
    }

    public getDotForDriverGUID(guid: string): JQuery<HTMLElement> | undefined {
        return this.dots.get(guid);
    }
}

const DriverGUIDDataKey = "driver-guid";

class LiveTimings implements WebsocketHandler {
    private readonly raceControl: RaceControl;
    private readonly liveMap: LiveMap;

    private readonly $connectedDriversTable: JQuery<HTMLTableElement>;
    private readonly $disconnectedDriversTable: JQuery<HTMLTableElement>;
    private readonly $storedTimes: JQuery<HTMLDivElement>;

    constructor(raceControl: RaceControl, liveMap: LiveMap) {
        this.raceControl = raceControl;
        this.liveMap = liveMap;
        this.$connectedDriversTable = $("#live-table");
        this.$disconnectedDriversTable = $("#live-table-disconnected");
        this.$storedTimes = $("#stored-times");

        setInterval(this.populateConnectedDrivers.bind(this), 1000);

        $(document).on("click", ".driver-link", this.toggleDriverSpeed.bind(this));
    }

    public handleWebsocketMessage(message: WSMessage): void {
        if (message.EventType === EventRaceControl) {
            this.populateConnectedDrivers();
            this.populateDisconnectedDrivers();
        }
    }

    private populateConnectedDrivers(): void {
        if (!this.raceControl.status || !this.raceControl.status.ConnectedDrivers) {
            return;
        }

        for (const driverGUID of this.raceControl.status.ConnectedDrivers.GUIDsInPositionalOrder) {
            const driver = this.raceControl.status.ConnectedDrivers.Drivers[driverGUID];

            if (!driver) {
                continue;
            }

            this.addDriverToTable(driver, this.$connectedDriversTable);
            this.populatePreviousLapsForDriver(driver);
        }

        if (this.raceControl.status.ConnectedDrivers.GUIDsInPositionalOrder.length > 0) {
            this.$connectedDriversTable.show();
        } else {
            this.$connectedDriversTable.hide();
        }
    }

    private populatePreviousLapsForDriver(driver: Driver): void {
        for (const carName in driver.Cars) {
            if (carName === driver.CarInfo.CarModel) {
                continue;
            }

            // create a fake new driver from the old driver. override details with their previous car
            // and add them to the disconnected drivers table. if the user rejoins in this car it will
            // be removed from the disconnected drivers table and placed into the connected drivers table.
            const dummyDriver = new Driver(driver);
            dummyDriver.CarInfo.CarModel = carName;

            this.addDriverToTable(dummyDriver, this.$disconnectedDriversTable);
        }
    }

    private populateDisconnectedDrivers(): void {
        if (!this.raceControl.status || !this.raceControl.status.DisconnectedDrivers) {
            return;
        }

        for (const driverGUID of this.raceControl.status.DisconnectedDrivers.GUIDsInPositionalOrder) {
            const driver = this.raceControl.status.DisconnectedDrivers.Drivers[driverGUID];

            if (!driver) {
                continue;
            }

            this.addDriverToTable(driver, this.$disconnectedDriversTable);
            this.populatePreviousLapsForDriver(driver);
        }

        if (this.$disconnectedDriversTable.find("tr").length > 1) {
            this.$storedTimes.show();
        } else {
            this.$storedTimes.hide();
        }
    }

    private addDriverToTable(driver: Driver, $table: JQuery<HTMLTableElement>): void {
        const addingDriverToDisconnectedTable = ($table === this.$disconnectedDriversTable);
        const driverID = driver.CarInfo.DriverGUID + "_" + driver.CarInfo.CarModel;
        const carInfo = driver.Cars[driver.CarInfo.CarModel];

        if (!carInfo) {
            return;
        }

        const $tr = $("<tr/>").attr({"id": driverID});

        // car position
        if (!addingDriverToDisconnectedTable) {
            const $tdPos = $("<td class='text-center'/>").text(driver.Position === 255 || driver.Position === 0 ? "" : driver.Position);
            $tr.append($tdPos);
        }

        // driver name
        const $tdName = $("<td/>").text(driver.CarInfo.DriverName);

        if (!addingDriverToDisconnectedTable) {
            // driver dot
            const driverDot = this.liveMap.getDotForDriverGUID(driver.CarInfo.DriverGUID);

            if (driverDot) {
                let dotClass = "dot";

                if (driverDot.find(".info").is(":hidden")) {
                    dotClass += " dot-inactive";
                }

                $tdName.prepend($("<div/>").attr({"class": dotClass}).css("background", randomColor({
                    luminosity: 'bright',
                    seed: driver.CarInfo.DriverGUID,
                })));
            }

            $tdName.attr("class", "driver-link");
            $tdName.data(DriverGUIDDataKey, driver.CarInfo.DriverGUID);
        }

        $tr.append($tdName);

        // car model
        const $tdCar = $("<td/>").text(prettifyName(driver.CarInfo.CarModel, true));
        $tr.append($tdCar);

        if (!addingDriverToDisconnectedTable) {
            let currentLapTimeText = "";

            if (moment(carInfo.LastLapCompletedTime).isAfter(moment(this.raceControl.status!.SessionStartTime))) {
                // only show current lap time text if the last lap completed time is after session start.
                currentLapTimeText = msToTime(moment().diff(moment(carInfo.LastLapCompletedTime)), false);
            }

            const $tdCurrentLapTime = $("<td/>").text(currentLapTimeText);
            $tr.append($tdCurrentLapTime);
        }

        if (!addingDriverToDisconnectedTable) {
            // last lap
            const $tdLastLap = $("<td/>").text(msToTime(carInfo.LastLap / 1000000));
            $tr.append($tdLastLap);
        }

        // best lap
        const $tdBestLapTime = $("<td/>").text(msToTime(carInfo.BestLap / 1000000));
        $tr.append($tdBestLapTime);

        if (!addingDriverToDisconnectedTable) {
            // gap
            const $tdGap = $("<td/>").text(driver.Split);
            $tr.append($tdGap);
        }

        // lap number
        const $tdLapNum = $("<td/>").text(carInfo.NumLaps ? carInfo.NumLaps : "0");
        $tr.append($tdLapNum);

        const $tdTopSpeedBestLap = $("<td/>").text(carInfo.TopSpeedBestLap ? carInfo.TopSpeedBestLap.toFixed(2) + "Km/h" : "");
        $tr.append($tdTopSpeedBestLap);

        if (!addingDriverToDisconnectedTable) {
            // events
            const $tdEvents = $("<td/>");

            if (moment(driver.LoadedTime).add("10", "seconds").isSameOrAfter(moment())) {
                // car just loaded
                let $tag = $("<span/>");
                $tag.attr({'class': 'badge badge-success live-badge'});
                $tag.text("Loaded");

                $tdEvents.append($tag);
            }

            if (driver.Collisions) {
                for (const collision of driver.Collisions) {
                    if (moment(collision.Time).add("10", "seconds").isSameOrAfter(moment())) {
                        let $tag = $("<span/>");
                        $tag.attr({'class': 'badge badge-danger live-badge'});
                        $tag.text(
                            "Crash " + collision.Type + " at " + collision.Speed.toFixed(2) + "Km/h"
                        );

                        $tdEvents.append($tag);
                    }
                }
            }

            $tr.append($tdEvents);
        }

        // remove any previous rows
        $("#" + driverID).remove();

        $table.append($tr);
    }

    private toggleDriverSpeed(e: ClickEvent): void {
        const $target = $(e.currentTarget);
        const driverGUID = $target.data(DriverGUIDDataKey);
        const $driverDot = this.liveMap.getDotForDriverGUID(driverGUID);

        if (!$driverDot) {
            return;
        }

        $driverDot.find(".info").toggle();
        $target.find(".dot").toggleClass("dot-inactive");
    }
}

function getAbbreviation(name: string): string {
    let parts = name.split(" ");

    if (parts.length < 1) {
        return name
    }

    let lastName = parts[parts.length - 1];

    return lastName.slice(0, 3).toUpperCase();
}

function randomColorForDriver(driverGUID: string): string {
    return randomColor({
        seed: driverGUID,
    })
}

const percentColors = [
    {pct: 0.25, color: {r: 0x00, g: 0x00, b: 0xff}},
    {pct: 0.625, color: {r: 0x00, g: 0xff, b: 0}},
    {pct: 1.0, color: {r: 0xff, g: 0x00, b: 0}}
];

function getColorForPercentage(pct: number) {
    let i;

    for (i = 1; i < percentColors.length - 1; i++) {
        if (pct < percentColors[i].pct) {
            break;
        }
    }

    let lower = percentColors[i - 1];
    let upper = percentColors[i];
    let range = upper.pct - lower.pct;
    let rangePct = (pct - lower.pct) / range;
    let pctLower = 1 - rangePct;
    let pctUpper = rangePct;
    let color = {
        r: Math.floor(lower.color.r * pctLower + upper.color.r * pctUpper),
        g: Math.floor(lower.color.g * pctLower + upper.color.g * pctUpper),
        b: Math.floor(lower.color.b * pctLower + upper.color.b * pctUpper)
    };

    return 'rgb(' + [color.r, color.g, color.b].join(',') + ')';
}