import { useContext } from "react";
import { gameContext } from "../contexts/gameContext";
import {
    isStateDenyAssassin,
    isStateDenyMoney,
    isStateDenySteal,
    isStateIdle,
    isStateStart,
} from "../utils/countState";

export const useAbility = () => {
    const { client } = useContext(gameContext);
    const clientState = client.state;
    const isIdle = isStateIdle(clientState);
    const isStart = isStateStart(clientState);
    const isDenyMoney = isStateDenyMoney(clientState);
    const isDenyAssassin = isStateDenyAssassin(clientState);
    const isDenySteal = isStateDenySteal(clientState);

    const isOthersRound = isIdle || !isStart;

    return {
        denyAssassin: isIdle || !isDenyAssassin,
        drawThreeCoins: isOthersRound,
        denyMoney: isIdle || !isDenyMoney,
        changeCard: isOthersRound,
        denySteal: isIdle || !isDenySteal,
        stealCoins: isOthersRound,
        assassin: isIdle || !isStart || client.coins < 3,
    };
};
