import {
    Button,
    FormControl,
    Input,
    MenuItem,
    Select,
    SelectChangeEvent,
} from "@mui/material";
import type { NextPage } from "next";
import styles from "../styles/Home.module.css";
import React, { useContext, useEffect, useState } from "react";
import { nakamaClient } from "../utils/nakama";
import { useRouter } from "next/router";
import { gameContext } from "../contexts/gameContext";

const Home: NextPage = () => {
    const [email, setEmail] = useState("");
    const [disableEmail, setDisableEmail] = useState(false);
    const [playerNum, setPlayerNum] = useState(2);
    const router = useRouter();
    const { counter, gameStart } = useContext(gameContext);

    useEffect(() => {
        const email = nakamaClient.getUserEmail();
        if (email) {
            setEmail(email);
            setDisableEmail(true);
        }
    }, []);

    useEffect(() => {
        if (gameStart) {
            router.push("in-game");
        }
    }, [gameStart]);

    const handlePlayerChange = (event: SelectChangeEvent<number>) => {
        setPlayerNum(event.target.value as number);
    };

    const handleChange = (
        e: React.ChangeEvent<HTMLTextAreaElement | HTMLInputElement>
    ) => {
        setEmail(e.target.value);
    };

    const handleClick = async () => {
        try {
            await nakamaClient.authenticate(email);
            await nakamaClient.findMatch(playerNum);
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
                        value={email}
                        placeholder="输入邮箱"
                        onChange={handleChange}
                        disabled={disableEmail}
                        error={!!email && !email.includes("@")}
                    />

                    <FormControl
                        variant="standard"
                        sx={{ m: 1, minWidth: 120 }}
                    >
                        <Select
                            labelId="demo-simple-select-standard-label"
                            id="demo-simple-select-standard"
                            value={playerNum}
                            onChange={handlePlayerChange}
                            label="Age"
                        >
                            <MenuItem value={2}>2人</MenuItem>
                            <MenuItem value={3}>3人</MenuItem>
                            <MenuItem value={4}>4人</MenuItem>
                        </Select>
                    </FormControl>

                    <Button
                        sx={{ width: "120px", m: 1 }}
                        variant="contained"
                        onClick={handleClick}
                        disabled={!email || !email.includes("@")}
                    >
                        加入比赛
                    </Button>
                </div>
            </main>
        </div>
    );
};

export default Home;
