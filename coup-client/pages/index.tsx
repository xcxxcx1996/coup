import { Button, Input } from "@mui/material";
import type { NextPage } from "next";
import styles from "../styles/Home.module.css";
import React, { useEffect, useState } from "react";
import { nakamaCLient } from "../utils/nakama";

const Home: NextPage = () => {
    const [name, setName] = useState("");
    const [disableUsername, setDisableUsername] = useState(false);
    useEffect(() => {
        const user = nakamaCLient.getUser();
        if (user) {
            setName(user.username as string);
            setDisableUsername(true);
        }
    }, []);
    const handleChange = (
        e: React.ChangeEvent<HTMLTextAreaElement | HTMLInputElement>
    ) => {
        setName(e.target.value);
    };
    const handleClick = async () => {
        try {
            await nakamaCLient.authenticate(name, !disableUsername);
            await nakamaCLient.findMatch();
        } catch (e) {
            console.log("-> e", e);
        }
    };
    return (
        <div className={styles.container}>
            <main className={styles.main}>
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
