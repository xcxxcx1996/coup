import { Client, Session, Socket } from "@heroiclabs/nakama-js";
import { v4 as uuidv4 } from "uuid";

class Nakama {
    private client: Client = new Client("defaultkey", "127.0.0.1", "7350");
    private session: Session;
    private socket: Socket;
    private matchID: any;
    private useSSL: boolean = false;
    constructor() {}

    async authenticate(username: string, create: boolean) {
        this.session = await this.client.authenticateCustom(
            "custom_id",
            create,
            username
        );
        localStorage.setItem("user_session", JSON.stringify(this.session));

        const trace = false;
        this.socket = this.client.createSocket(this.useSSL, trace);
        await this.socket.connect(this.session, false);
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

    async makeMove(index: number) {
        // ep4
        var data = { position: index };
        await this.socket.sendMatchState(this.matchID, 4, data);
        console.log("Match data sent");
    }
}

export const nakamaCLient = new Nakama();
