import {Player} from "./Player";

export interface Team {
    id: number;
    name: string;
    tag: string;
    description: string;
    wins: number;
    losses: number;
    iconSmall: string;
    iconLarge: string;
    players: Player[];
    substitutes: Player[];
    visible: boolean;
}

export function EmptyTeam(): Team {
    return {
        id: 0,
        name: "",
        tag: "",
        description: "",
        wins: 0,
        losses: 0,
        iconSmall: "",
        iconLarge: "",
        players: [],
        substitutes: [],
        visible: false,
    }
}
