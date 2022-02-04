import { Box, Button } from "@mui/material";
import React, { useState } from "react";
import { nakamaClient } from "../../utils/nakama";
export interface GameHistoryProps {}

export const GameHistory: React.FC<GameHistoryProps> = ({}) => {
    const [gameStart, setGameStart] = useState(false);
    const handleClick = async () => {
        setGameStart(true);
        await nakamaClient.startGame();
    };
    return (
        <Box
            sx={{
                m: 1,
                p: 1,
                display: "flex",
                flexDirection: "column",
                border: "1px solid",
                borderRadius: 2,
                height: "300px",
                overflow: "scroll",
            }}
        >
            {!gameStart && (
                <Button variant="contained" onClick={handleClick}>
                    开始游戏
                </Button>
            )}
        </Box>
    );
};
