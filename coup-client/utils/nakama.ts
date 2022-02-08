import { Client, Session, Socket } from "@heroiclabs/nakama-js";
import { OP_CODE } from "../constants/op_code";

const saveInStorage = (key: string, value: string): void => {
    localStorage.setItem(key, value);
};

const retrieveInStorage = (key: string): string => {
    return localStorage.getItem(key);
};

class Nakama {
    private client: Client = new Client("defaultkey", "localhost", "7350");
    public session: Session;
    private useSSL: boolean = false;
    public socket: Socket = this.client.createSocket(this.useSSL, false);
    private matchID: any;
    private password: string = "password123";
    constructor() {}

    authenticate = async (email: string) => {
        this.session = await this.client.authenticateEmail(
            email,
            this.password,
            true,
            email.split("@").shift()
        );
        const token = this.session.token;
        const refreshToken = this.session.refresh_token;
        saveInStorage("email", email);
        saveInStorage("token", token);
        saveInStorage("refreshToken", refreshToken);
        await this.socket.connect(this.session, true);
    };

    reconnect = async () => {
        const email = retrieveInStorage("email");
        if (email) {
            const token = retrieveInStorage("token");
            const refreshToken = retrieveInStorage("refreshToken");
            this.session = Session.restore(token, refreshToken);
            const matchID = retrieveInStorage("matchID");
            this.socket
                .connect(this.session, true)
                .then(() => {
                    this.socket
                        .joinMatch(matchID)
                        .then((result) => {
                            console.log("-> result", result);
                        })
                        .catch((err) => console.log("join failed", err));
                })
                .catch((err) => console.log(err));
        }
    };

    getUserEmail() {
        return retrieveInStorage("email");
    }

    findMatch = async () => {
        // ep4
        const rpcid = "find_match";
        const matches = await this.client.rpc(this.session, rpcid, {});
        this.matchID = (matches.payload as { matchIds: string[] }).matchIds[0];
        saveInStorage("matchID", this.matchID);
        await this.socket.joinMatch(this.matchID);
    };

    readyStart = async () => {
        await this.socket.sendMatchState(
            this.matchID,
            OP_CODE.READY_START,
            null
        );
    };

    startGame = async () => {
        await this.socket.sendMatchState(this.matchID, OP_CODE.START, null);
    };

    async rejected() {
        await this.socket.sendMatchState(this.matchID, OP_CODE.REJECTED, null);
    }

    assassin = async (playerId: string) => {
        await this.socket.sendMatchState(
            this.matchID,
            OP_CODE.ASSASSIN,
            playerId
        );
    };

    denyKill = async () => {
        await this.socket.sendMatchState(this.matchID, OP_CODE.DENY_KILL, null);
    };

    drawThreeCoins = async () => {
        await this.socket.sendMatchState(
            this.matchID,
            OP_CODE.DRAW_THREE_COINS,
            null
        );
    };

    denyMoney = async () => {
        await this.socket.sendMatchState(
            this.matchID,
            OP_CODE.DENY_MONEY,
            null
        );
    };

    changeCard = async () => {
        await this.socket.sendMatchState(
            this.matchID,
            OP_CODE.CHANGE_CARD,
            null
        );
    };

    denySteal = async () => {
        await this.socket.sendMatchState(
            this.matchID,
            OP_CODE.DENY_STEAL,
            null
        );
    };

    stealCoins = async (playerId: string) => {
        await this.socket.sendMatchState(
            this.matchID,
            OP_CODE.STEAL_COINS,
            playerId
        );
    };

    discardCard = async (cardId: string) => {
        await this.socket.sendMatchState(
            this.matchID,
            OP_CODE.DISCARD_CARD,
            cardId
        );
    };

    drawCoins = async (coinNum: number) => {
        await this.socket.sendMatchState(
            this.matchID,
            OP_CODE.DRAW_COINS,
            coinNum
        );
    };

    question = async (isQuestion: boolean) => {
        await this.socket.sendMatchState(
            this.matchID,
            OP_CODE.QUESTION,
            isQuestion
        );
    };

    coup = async (playerId: string) => {
        await this.socket.sendMatchState(this.matchID, OP_CODE.COUP, playerId);
    };
}

export const nakamaClient = new Nakama();
