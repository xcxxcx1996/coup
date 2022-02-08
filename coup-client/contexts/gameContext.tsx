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
import { ROLES, rolesMap } from "../constants";

export interface GameContext {
    users: IUser[];
    setUsers: Dispatch<SetStateAction<IUser[]>>;
    currentPlayer: string;
    timeLeft: number;
    shouldReconnect: boolean;
    chooseCards: ICard[];
    infos: string[];
    client: PlayerInfo;
}

export const gameContext = createContext<GameContext>({
    users: null,
    setUsers: null,
    currentPlayer: null,
    timeLeft: null,
    shouldReconnect: null,
    chooseCards: null,
    infos: [],
    client: null,
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
    const [currentPlayer, setCurrentPlayer] = useState("");
    const [timeLeft, setTimeLeft] = useState(0);
    const [shouldReconnect, setShouldReconnect] = useState(false);
    const [chooseCards, setChooseCards] = useState<ICard[]>([]);
    const [infos, setInfos] = useState<string[]>([]);
    const [client, setClient] = useState<PlayerInfo>({
        id: "",
        coins: 0,
        cards: [],
        name: "",
        state: 0,
    });
    const userId = nakamaClient?.session?.user_id;
    useEffect(() => {
        nakamaClient.socket.onmatchdata = (matchData: MatchData) => {
            console.log("-> matchData", matchData);
            switch (matchData.op_code) {
                case OP_CODE.START:
                case OP_CODE.UPDATE:
                    const playerInfos: PlayerInfos = matchData.data.playerInfos;
                    const currentPlayerId: string =
                        matchData.data.currentPlayerId;
                    const currentPlayer = playerInfos[currentPlayerId];
                    const clientPlayer = playerInfos[userId];
                    setClient(clientPlayer);
                    setUsers(transformPlayerInfos(playerInfos, userId));
                    setCurrentPlayer(currentPlayer.name);
                    break;
                case OP_CODE.TICK:
                    setTimeLeft(matchData.data.deadline);
                    break;
                case OP_CODE.CHANGE_CARD:
                    const chooseCards = matchData.data.chooseCards;
                    setChooseCards(chooseCards);
                    break;
                case OP_CODE.INFO:
                    setInfos((infos) => [matchData.data.message, ...infos]);
                    break;
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
                client,
                chooseCards,
                infos,
            }}
        >
            {children}
        </gameContext.Provider>
    );
};
