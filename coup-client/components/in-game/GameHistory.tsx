import { Box, Button, Divider, Theme } from "@mui/material";
import React, { useCallback, useContext } from "react";
import { gameContext } from "../../contexts/gameContext";
import { nakamaClient } from "../../utils/nakama";
import { styled } from "@mui/material/styles";

const Root = styled("div")(({ theme }: { theme: Theme }) => ({
    width: "100%",
    ...theme.typography.body2,
    "& > :not(style) + :not(style)": {
        marginTop: theme.spacing(2),
    },
}));

export const GameHistory = () => {
    const { shouldReconnect, infos } = useContext(gameContext);
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
            <Root>
                {infos.map((info, index) => (
                    <div key={index}>
                        {info}
                        <Divider />
                    </div>
                ))}
            </Root>
        </Box>
    );
};
