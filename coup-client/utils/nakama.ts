import { Client, Session, Socket } from "@heroiclabs/nakama-js";
import { v4 as uuidv4 } from "uuid";
import { OP_CODE } from "../constants/op_code";
import { MatchData } from "@heroiclabs/nakama-js/socket";

class Nakama {
    private client: Client = new Client("defaultkey", "localhost", "7350");
    private session: Session;
    private useSSL: boolean = false;
    public socket: Socket = this.client.createSocket(this.useSSL, false);
    private matchID: any;
    constructor() {}

    async authenticate(username: string) {
        let deviceId = localStorage.getItem("deviceId");
        let create = false;
        if (!deviceId) {
            deviceId = uuidv4();
            localStorage.setItem("deviceId", deviceId);
            create = true;
        }
        this.session = await this.client.authenticateDevice(
            deviceId,
            create,
            username
        );
        localStorage.setItem("user_session", JSON.stringify(this.session));
        await this.socket.connect(this.session, true);
    }

    getUser() {
        const sessionString = localStorage.getItem("user_session");
        return sessionString ? JSON.parse(sessionString) : null;
    }

    async findMatch() {
        // ep4
        const rpcid = "find_match";
        const matches = await this.client.rpc(this.session, rpcid, {});
        this.matchID = (matches.payload as { matchIds: string[] }).matchIds[0];
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

    async assassin() {
        await this.socket.sendMatchState(this.matchID, OP_CODE.ASSASSIN, null);
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

    async stealCoins() {
        await this.socket.sendMatchState(
            this.matchID,
            OP_CODE.STEAL_COINS,
            null
        );
    }

    async discardCard() {
        await this.socket.sendMatchState(
            this.matchID,
            OP_CODE.DISCARD_CARD,
            null
        );
    }

    async drawCoins() {
        await this.socket.sendMatchState(
            this.matchID,
            OP_CODE.DRAW_COINS,
            null
        );
    }

    async question() {
        await this.socket.sendMatchState(this.matchID, OP_CODE.QUESTION, null);
    }

    async finishGame() {
        await this.socket.sendMatchState(this.matchID, OP_CODE.DONE, null);
    }

    async coup() {
        await this.socket.sendMatchState(this.matchID, OP_CODE.COUP, null);
    }
}

export const nakamaClient = new Nakama();
