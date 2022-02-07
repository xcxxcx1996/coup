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

    async authenticate(email: string) {
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
    }

    reconnect() {
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
    }

    getUserEmail() {
        return retrieveInStorage("email");
    }

    async findMatch() {
        // ep4
        const rpcid = "find_match";
        const matches = await this.client.rpc(this.session, rpcid, {});
        this.matchID = (matches.payload as { matchIds: string[] }).matchIds[0];
        saveInStorage("matchID", this.matchID);
        await this.socket.joinMatch(this.matchID);
    }

    async readyStart() {
        await this.socket.sendMatchState(
            this.matchID,
            OP_CODE.READY_START,
            null
        );
    }

    async startGame() {
        await this.socket.sendMatchState(this.matchID, OP_CODE.START, null);
    }

    async update() {
        await this.socket.sendMatchState(this.matchID, OP_CODE.UPDATE, null);
    }

    async rejected() {
        await this.socket.sendMatchState(this.matchID, OP_CODE.REJECTED, null);
    }

    async drawCard() {
        await this.socket.sendMatchState(this.matchID, OP_CODE.DRAW_CARD, null);
    }

    async assassin(playerId: string) {
        await this.socket.sendMatchState(
            this.matchID,
            OP_CODE.ASSASSIN,
            playerId
        );
    }

    async denyKill() {
        await this.socket.sendMatchState(this.matchID, OP_CODE.DENY_KILL, null);
    }

    async drawThreeCoins() {
        await this.socket.sendMatchState(
            this.matchID,
            OP_CODE.DRAW_THREE_COINS,
            null
        );
    }

    async denyMoney() {
        await this.socket.sendMatchState(
            this.matchID,
            OP_CODE.DENY_MONEY,
            null
        );
    }

    async changeCard() {
        await this.socket.sendMatchState(
            this.matchID,
            OP_CODE.CHANGE_CARD,
            null
        );
    }

    async denySteal() {
        await this.socket.sendMatchState(
            this.matchID,
            OP_CODE.DENY_STEAL,
            null
        );
    }

    async stealCoins(playerId: string) {
        await this.socket.sendMatchState(
            this.matchID,
            OP_CODE.STEAL_COINS,
            playerId
        );
    }

    async discardCard(cardId: string) {
        await this.socket.sendMatchState(
            this.matchID,
            OP_CODE.DISCARD_CARD,
            cardId
        );
    }

    async drawCoins(coinNum: number) {
        await this.socket.sendMatchState(
            this.matchID,
            OP_CODE.DRAW_COINS,
            coinNum
        );
    }

    async question(isQuestion: boolean) {
        await this.socket.sendMatchState(
            this.matchID,
            OP_CODE.QUESTION,
            isQuestion
        );
    }

    async finishGame() {
        await this.socket.sendMatchState(this.matchID, OP_CODE.DONE, null);
    }

    async coup(playerId: string) {
        await this.socket.sendMatchState(this.matchID, OP_CODE.COUP, playerId);
    }

    async info() {
        await this.socket.sendMatchState(this.matchID, OP_CODE.INFO, null);
    }
}

export const nakamaClient = new Nakama();
