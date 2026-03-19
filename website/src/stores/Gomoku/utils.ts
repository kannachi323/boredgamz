import { LobbyRequest, Player } from "@/pages/Games/Gomoku/types";

const UNITS_IN_SECONDS: {[key: string]: number} = {
    'nanoseconds': 1e-9,
    'microseconds': 1e-6,
    'milliseconds': 1e-3,
    'seconds': 1, 
    'minutes': 60,
    'hours': 3600,
};

export function convertTime(val: number, unit1: string, unit2: string) : number {
    const unit1Lower = unit1.toLowerCase();
    const unit2Lower = unit2.toLowerCase();
    
    const factor1 = UNITS_IN_SECONDS[unit1Lower];
    const factor2 = UNITS_IN_SECONDS[unit2Lower];

    const valInSeconds = val * factor1;

    const convertedVal = valInSeconds / factor2;

    return convertedVal;
}

export function createPlayer() : Player {
    return {
        playerID: '',
        playerName: '',
        color: 'black',
        playerClock: { remaining: 0 },
    }
}

export function createLobbyRequest() : LobbyRequest {
    return {
        type: "lobby",
        data: {
            playerName: "",
            playerID: "",
            playerColor: "black",
            mode: "casual",
            timeControl: "Rapid",
            name: "9x9",
        }
    }
}