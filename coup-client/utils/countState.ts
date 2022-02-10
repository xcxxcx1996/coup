import { State } from "../constants/state";

export const isStateIdle = (playerState: number): boolean => {
    return playerState === State.IDLE;
};
export const isStateQuestion = (playerState: number): boolean => {
    return playerState === State.QUESTION;
};
export const isStateDenyMoney = (playerState: number): boolean => {
    return playerState === State.DENY_MONEY;
};
export const isStateDenyAssassin = (playerState: number): boolean => {
    return playerState === State.DENY_ASSASSIN;
};
export const isStateDenySteal = (playerState: number): boolean => {
    return playerState === State.DENY_STEAL;
};
export const isStateDiscard = (playerState: number): boolean => {
    return playerState === State.DISCARD;
};
export const isStateStart = (playerState: number): boolean => {
    return playerState === State.START;
};
export const isStateChooseCard = (playerState: number): boolean => {
    return playerState === State.CHOOSE_CARD;
};
