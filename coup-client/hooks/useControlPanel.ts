import { useContext } from "react";
import { gameContext } from "../contexts/gameContext";
import {
    isStateDiscard,
    isStateIdle,
    isStateQuestion,
    isStateStart,
} from "../utils/countState";

export const useControlPanel = () => {
    const { client } = useContext(gameContext);
    const { state } = client;
    const isIdle = isStateIdle(state);
    const isQuestion = isStateQuestion(state);
    const isDiscard = isStateDiscard(state);
    const isStart = isStateStart(state);
    const isMustCoup = client.coins >= 10 && isStart;

    const isOthersRound = isIdle || !isStart;

    return {
        drawCoins: isOthersRound || isMustCoup,
        useAbility: isIdle || isQuestion || isDiscard || isMustCoup,
        question: isIdle || !isQuestion,
        coup: isOthersRound || client.coins < 7,
        discard: isIdle || !isDiscard,
    };
};
