import {
    createContext,
    Dispatch,
    FC,
    SetStateAction,
    useEffect,
    useState,
} from "react";
import { IUser } from "../components/in-game/UserCarousel";
import { nakamaClient } from "../utils/nakama";
import { MatchData } from "@heroiclabs/nakama-js/socket";
import { OP_CODE } from "../constants/op_code";
import { rolesMap } from "../constants";

export interface GameContext {
    users: IUser[];
    setUsers: Dispatch<SetStateAction<IUser[]>>;
    currentPlayer: string;
    timeLeft: number;
    shouldReconnect: boolean;
    cards: ICard[];
    setCards: Dispatch<SetStateAction<ICard[]>>;
    shouldDiscard: boolean;
    chooseCards: ICard[];
}

export const gameContext = createContext<GameContext>({
    users: null,
    setUsers: null,
    currentPlayer: null,
    timeLeft: null,
    shouldReconnect: null,
    cards: null,
    setCards: null,
    shouldDiscard: false,
    chooseCards: null,
});

export interface PlayerInfo {
    [id: string]: {
        id: string;
        coins: number;
        cards: ICard[];
        name: string;
    };
}

export interface ICard {
    id: string;
    role: number;
}

export const transformPlayerInfos = (playerInfo: PlayerInfo) => {
    return Object.values(playerInfo).map((item) => ({
        id: item.id,
        name: item.name,
        coins: item.coins,
        roles: item.cards.map((card) => rolesMap[String(card.role)]),
    }));
};

export const GameContextProvider: FC = ({ children }) => {
    const [users, setUsers] = useState<IUser[]>([]);
    const [cards, setCards] = useState<ICard[]>([]);
    const [currentPlayer, setCurrentPlayer] = useState("");
    const [timeLeft, setTimeLeft] = useState(0);
    const [shouldReconnect, setShouldReconnect] = useState(false);
    const [shouldDiscard, setShouldDiscard] = useState(false);
    const [chooseCards, setChooseCards] = useState<ICard[]>([]);
    useEffect(() => {
        nakamaClient.socket.onmatchdata = (matchData: MatchData) => {
            console.log("-> matchData", matchData);
            switch (matchData.op_code) {
                case OP_CODE.START:
                case OP_CODE.UPDATE:
                    const playerInfos: PlayerInfo = matchData.data.playerInfos;
                    const currentPlayerId: string =
                        matchData.data.currentPlayerId;
                    setUsers(transformPlayerInfos(playerInfos));
                    setCards(playerInfos[currentPlayerId].cards);
                    setCurrentPlayer(playerInfos[currentPlayerId].name);
                    break;
                case OP_CODE.TICK:
                    setTimeLeft(matchData.data.deadline);
                    break;
                case OP_CODE.DISCARD_CARD:
                    setShouldDiscard(true);
                    break;
                case OP_CODE.CHANGE_CARD:
                    const chooseCards = matchData.data.chooseCards;
                    setChooseCards(chooseCards);
            }
        };
    }, []);

    useEffect(() => {
        if (!nakamaClient.session) {
            setShouldReconnect(true);
        } else {
            setShouldReconnect(false);
        }
    }, [nakamaClient.session, setShouldReconnect]);

    return (
        <gameContext.Provider
            value={{
                users,
                setUsers,
                currentPlayer,
                timeLeft,
                shouldReconnect,
                cards,
                setCards,
                shouldDiscard,
                chooseCards,
            }}
        >
            {children}
        </gameContext.Provider>
    );
};
