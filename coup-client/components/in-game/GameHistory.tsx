import { Box, Button } from "@mui/material";
import React, { useCallback, useContext } from "react";
import { gameContext } from "../../contexts/gameContext";
import { nakamaClient } from "../../utils/nakama";

export const GameHistory = () => {
    const { shouldReconnect } = useContext(gameContext);
    const handleReconnect = useCallback(() => {
        nakamaClient.reconnect();
    }, []);
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
            {shouldReconnect && (
                <Button
                    variant="contained"
                    color={"error"}
                    onClick={handleReconnect}
                >
                    重连
                </Button>
            )}
            Game start!
        </Box>
    );
};
