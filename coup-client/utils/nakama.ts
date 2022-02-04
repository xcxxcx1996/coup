import { Client, Session, Socket } from "@heroiclabs/nakama-js";
import { v4 as uuidv4 } from "uuid";
import { OP_CODE } from "../constants/op_code";

class Nakama {
    private client: Client = new Client("defaultkey", "localhost", "7350");
    private session: Session;
    private socket: Socket;
    private matchID: any;
    private useSSL: boolean = false;
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

        const trace = false;
        this.socket = this.client.createSocket(this.useSSL, trace);
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
        console.log("Matched joined!");
    }

    async startGame() {
        await this.socket.sendMatchState(this.matchID, OP_CODE.START, null);
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

    async drawTwoCoins() {
        await this.socket.sendMatchState(
            this.matchID,
            OP_CODE.DRAW_TWO_COINS,
            null
        );
    }

    async denyTwoCoins() {
        await this.socket.sendMatchState(
            this.matchID,
            OP_CODE.DENY_TWO_COINS,
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

    async drawCoin() {
        await this.socket.sendMatchState(this.matchID, OP_CODE.DRAW_COIN, null);
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
