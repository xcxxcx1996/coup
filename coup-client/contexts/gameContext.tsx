import {
    createContext,
    Dispatch,
    FC,
    SetStateAction,
    useEffect,
    useState,
} from "react";
import { IUser } from "../components/in-game/UserCarousel";
import { nakamaClient, retrieveInStorage } from "../utils/nakama";
import { MatchData } from "@heroiclabs/nakama-js/socket";
import { OP_CODE } from "../constants/op_code";
import { ROLES, rolesMap } from "../constants";

export interface GameContext {
    users: IUser[];
    currentPlayer: PlayerInfo;
    timeLeft: number;
    shouldReconnect: boolean;
    setShouldReconnect: Dispatch<SetStateAction<boolean>>;
    chooseCards: ICard[];
    infos: string[];
    client: PlayerInfo;
    gameEnd: boolean;
    gameStart: boolean;
    counter: number;
}

const initialPlayer: PlayerInfo = {
    id: "",
    coins: 0,
    cards: [],
    name: "",
    state: 0,
};

export const gameContext = createContext<GameContext>({
    users: null,
    currentPlayer: null,
    timeLeft: null,
    shouldReconnect: null,
    setShouldReconnect: null,
    chooseCards: null,
    infos: [],
    client: null,
    gameEnd: false,
    gameStart: false,
    counter: null,
});

export interface PlayerInfo {
    id: string;
    coins: number;
    cards: ICard[];
    name: string;
    state?: number;
}

export interface PlayerInfos {
    [id: string]: PlayerInfo;
}

export interface ICard {
    id: string;
    role: number;
}

export const transformPlayerInfos = (
    playerInfo: PlayerInfos,
    userId: string
) => {
    return Object.values(playerInfo).map((item) => ({
        id: item.id,
        name: item.name,
        coins: item.coins,
        roles: item.cards.map((card) => {
            return item.id === userId
                ? rolesMap[String(card.role)]
                : rolesMap[ROLES.UNROLE];
        }),
    }));
};

export const GameContextProvider: FC = ({ children }) => {
    const [users, setUsers] = useState<IUser[]>([]);
    const [currentPlayer, setCurrentPlayer] =
        useState<PlayerInfo>(initialPlayer);
    const [timeLeft, setTimeLeft] = useState(0);
    const [shouldReconnect, setShouldReconnect] = useState(false);
    const [chooseCards, setChooseCards] = useState<ICard[]>([]);
    const [infos, setInfos] = useState<string[]>([]);
    const [client, setClient] = useState<PlayerInfo>(initialPlayer);
    const [gameStart, setGameStart] = useState(false);
    const [gameEnd, setGameEnd] = useState(false);
    const [counter, setCounter] = useState<number>(null);
    useEffect(() => {
        nakamaClient.socket.onmatchdata = (matchData: MatchData) => {
            console.log("-> matchData", matchData);
            switch (matchData.op_code) {
                case OP_CODE.READY_START:
                    const nextGameStart = matchData.data.nextGameStart;
                    setCounter(parseInt(nextGameStart));
                    break;
                case OP_CODE.START:
                case OP_CODE.UPDATE:
                    if (matchData.op_code === OP_CODE.START) {
                        setGameStart(true);
                        setCounter(null);
                    }
                    const playerInfos: PlayerInfos = matchData.data.playerInfos;
                    const userId = nakamaClient.session
                        ? nakamaClient.session.user_id
                        : retrieveInStorage("userId");
                    const currentPlayerId: string =
                        matchData.data.currentPlayerId;
                    const currentPlayer = playerInfos[currentPlayerId];
                    const clientPlayer = playerInfos[userId];
                    setClient(clientPlayer);
                    setUsers(transformPlayerInfos(playerInfos, userId));
                    setCurrentPlayer(currentPlayer);
                    break;
                case OP_CODE.TICK:
                    setTimeLeft(matchData.data.deadline);
                    break;
                case OP_CODE.CHOOSE_CARD:
                    const chooseCards = matchData.data.cards;
                    setChooseCards(chooseCards || []);
                    break;
                case OP_CODE.INFO:
                    setInfos((infos) => [...infos, matchData.data.message]);
                    break;
                case OP_CODE.DONE:
                    const winner = matchData.data.winner;
                    setInfos((infos) => [
                        ...infos,
                        `${winner.name}???????????????????????????`,
                    ]);
                    setGameStart(false);
                    setGameEnd(true);
            }
        };
        return () => {
            nakamaClient.socket.disconnect(true);
        };
    }, []);

    useEffect(() => {
        if (!nakamaClient.session) {
            setShouldReconnect(true);
        } else {
            setShouldReconnect(false);
        }
    }, [nakamaClient.session, setShouldReconnect]);

    useEffect(() => {
        nakamaClient.socket.ondisconnect = (evt: Event) => {
            if (evt.type === "close") {
                setShouldReconnect(true);
            }
        };
    }, []);

    return (
        <gameContext.Provider
            value={{
                users,
                currentPlayer,
                timeLeft,
                shouldReconnect,
                setShouldReconnect,
                client,
                chooseCards,
                infos,
                gameEnd,
                gameStart,
                counter,
            }}
        >
            {children}
        </gameContext.Provider>
    );
};
