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
