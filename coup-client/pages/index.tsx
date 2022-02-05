import { Button, Input } from "@mui/material";
import type { NextPage } from "next";
import styles from "../styles/Home.module.css";
import React, { useCallback, useEffect, useState } from "react";
import { nakamaClient } from "../utils/nakama";
import { useRouter } from "next/router";
import { MatchData } from "@heroiclabs/nakama-js/socket";
import { OP_CODE } from "../constants/op_code";

const Home: NextPage = () => {
    const [name, setName] = useState("");
    const [disableUsername, setDisableUsername] = useState(false);
    const [counter, setCounter] = useState(0);
    const router = useRouter();
    const matchDataHandler = useCallback(
        (matchData: MatchData) => {
            console.log("-> matchData111", matchData);
            if (matchData.op_code === OP_CODE.READY_START) {
                const nextGameStart = matchData.data.nextGameStart;
                setCounter(Math.floor(nextGameStart - Date.now() / 1000));
            }
        },
        [setCounter]
    );
    useEffect(() => {
        nakamaClient.socket.onmatchdata = matchDataHandler;
        return () => {
            setCounter(0);
        };
    }, [matchDataHandler, setCounter]);

    useEffect(() => {
        const user = nakamaClient.getUser();
        if (user) {
            setName(user.username as string);
            setDisableUsername(true);
        }
    }, []);

    useEffect(() => {
        if (counter < 0) {
            router.push("in-game").then((r) => {
                console.log("-> r", r);
            });
        }
    }, [counter, setCounter]);

    const handleChange = (
        e: React.ChangeEvent<HTMLTextAreaElement | HTMLInputElement>
    ) => {
        setName(e.target.value);
    };
    const handleClick = async () => {
        try {
            await nakamaClient.authenticate(name);
            await nakamaClient.findMatch();
            // await router.push("in-game");
        } catch (e) {
            console.log("-> e", e);
        }
    };
    return (
        <div className={styles.container}>
            <main className={styles.main}>
                {counter >= 0 && <h2>{counter}</h2>}
                <h1 className={styles.title}>政变</h1>
                <div className={styles.grid}>
                    <Input
                        sx={{ mt: 3 }}
                        value={name}
                        placeholder="输入用户名"
                        onChange={handleChange}
                        disabled={disableUsername}
                    />
                    <Button
                        sx={{ width: "120px", m: 1 }}
                        variant="contained"
                        onClick={handleClick}
                        disabled={!name}
                    >
                        加入比赛
                    </Button>
                </div>
            </main>
        </div>
    );
};

export default Home;
