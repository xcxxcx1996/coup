import { Box } from "@mui/material";
import React from "react";
export interface GameHistoryProps {}

export const GameHistory: React.FC<GameHistoryProps> = ({}) => {
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
            Game start!
        </Box>
    );
};
