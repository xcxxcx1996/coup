import { Client, Session, Socket } from "@heroiclabs/nakama-js";
import { OP_CODE } from "../constants/op_code";

export const saveInStorage = (key: string, value: string): void => {
    localStorage.setItem(key, value);
};

export const retrieveInStorage = (key: string): string => {
    return localStorage.getItem(key);
};

const HOST = process.env.NEXT_PUBLIC_HOST || "localhost";

class Nakama {
    private client: Client = new Client("defaultkey", HOST, "7350");
    public session: Session;
    private useSSL: boolean = false;
    public socket: Socket = this.client.createSocket(this.useSSL, false);
    private matchID: any;
    private password: string = "password123";
    constructor() {}

    authenticate = async (email: string) => {
        const isCreate = email !== retrieveInStorage("email");
        this.session = await this.client.authenticateEmail(
            email,
            this.password,
            isCreate,
            email.split("@").shift()
        );
        const token = this.session.token;
        const refreshToken = this.session.refresh_token;
        const userId = this.session.user_id;
        saveInStorage("userId", userId);
        saveInStorage("email", email);
        saveInStorage("token", token);
        saveInStorage("refreshToken", refreshToken);
    };

    restoreSessionOrAuthenticate = async () => {
        const email = retrieveInStorage("email");
        await this.authenticate(email);
    };

    reconnect = async () => {
        const matchID = retrieveInStorage("matchID");
        this.matchID = matchID;
        await this.socket.connect(this.session, true);
        await this.socket.joinMatch(matchID);
    };

    getUserEmail() {
        return retrieveInStorage("email");
    }

    findMatch = async (playerNum: number) => {
        const rpcid = "find_match";
        const matches = await this.client.rpc(this.session, rpcid, {
            player_num: playerNum,
        });
        console.log("matches:",matches)
        this.matchID = (matches.payload as { matchIds: string[] }).matchIds[0];
        saveInStorage("matchID", this.matchID);
        await this.socket.connect(this.session, true);
        await this.socket.joinMatch(this.matchID);
    };

    leaveMatch = async () => {
        await this.socket.leaveMatch(this.matchID);
    };

    assassin = async (playerId: string) => {
        await this.socket.sendMatchState(this.matchID, OP_CODE.ASSASSIN, {
            player_id: playerId,
        });
    };

    denyKill = async (isDeny: boolean) => {
        await this.socket.sendMatchState(this.matchID, OP_CODE.DENY_KILL, {
            is_deny: isDeny,
        });
    };

    drawThreeCoins = async () => {
        await this.socket.sendMatchState(
            this.matchID,
            OP_CODE.DRAW_THREE_COINS,
            null
        );
    };

    denyMoney = async (isDeny: boolean) => {
        await this.socket.sendMatchState(this.matchID, OP_CODE.DENY_MONEY, {
            is_deny: isDeny,
        });
    };

    changeCard = async () => {
        await this.socket.sendMatchState(
            this.matchID,
            OP_CODE.CHANGE_CARD,
            null
        );
    };

    chooseCard = async (cards: string[]) => {
        await this.socket.sendMatchState(this.matchID, OP_CODE.CHOOSE_CARD, {
            cards,
        });
    };

    denySteal = async (isDeny: boolean, role: string) => {
        await this.socket.sendMatchState(this.matchID, OP_CODE.DENY_STEAL, {
            is_deny: isDeny,
            role,
        });
    };

    stealCoins = async (playerId: string) => {
        await this.socket.sendMatchState(this.matchID, OP_CODE.STEAL_COINS, {
            player_id: playerId,
        });
    };

    discardCard = async (cardId: string) => {
        await this.socket.sendMatchState(this.matchID, OP_CODE.DISCARD_CARD, {
            card_id: cardId,
        });
    };

    drawCoins = async (coinNum: number) => {
        await this.socket.sendMatchState(this.matchID, OP_CODE.DRAW_COINS, {
            coins: coinNum,
        });
    };

    question = async (isQuestion: boolean) => {
        await this.socket.sendMatchState(this.matchID, OP_CODE.QUESTION, {
            is_question: isQuestion,
        });
    };

    coup = async (playerId: string) => {
        await this.socket.sendMatchState(this.matchID, OP_CODE.COUP, {
            player_id: playerId,
        });
    };
}

export const nakamaClient = new Nakama();
