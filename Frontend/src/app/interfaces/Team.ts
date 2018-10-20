import {Player} from "./Player";

export interface Team {
    id: number;
    name: string;
    tag: string;
    wins: number;
    losses: number;
    players: Player[];
    substitutes: Player[];
    visible: boolean;
}
