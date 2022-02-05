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
}

export const gameContext = createContext<GameContext>({
    users: null,
    setUsers: null,
    currentPlayer: null,
    timeLeft: null,
});

export interface PlayerInfo {
    [id: string]: {
        id: string;
        coins: number;
        cards: { id: string; role: number }[];
    };
}

export const GameContextProvider: FC = ({ children }) => {
    const [users, setUsers] = useState<IUser[]>([]);
    const [currentPlayer, setCurrentPlayer] = useState("");
    const [timeLeft, setTimeLeft] = useState(0);
    useEffect(() => {
        nakamaClient.socket.onmatchdata = (matchData: MatchData) => {
            if (matchData.op_code === OP_CODE.UPDATE) {
                console.log("-> matchData", matchData);
                const playerInfos: PlayerInfo = matchData.data.playerInfos;
                const currentPlayerId = matchData.data.currentPlayerId;
                const deadline = matchData.data.deadline;
                const users: IUser[] = Object.values(playerInfos).map(
                    (item) => ({
                        name: item.id.slice(0, 5),
                        coins: item.coins,
                        roles: item.cards.map(
                            (card) => rolesMap[String(card.role)]
                        ),
                    })
                );
                setUsers(users);
                setCurrentPlayer(currentPlayerId.slice(0, 5));
                setTimeLeft(deadline - Math.floor(Date.now() / 1000));
            }
        };
    }, []);

    return (
        <gameContext.Provider
            value={{ users, setUsers, currentPlayer, timeLeft }}
        >
            {children}
        </gameContext.Provider>
    );
};
