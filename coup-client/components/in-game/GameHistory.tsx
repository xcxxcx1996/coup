import { Box, Button, Divider, Theme } from "@mui/material";
import React, {
    useCallback,
    useContext,
    useEffect,
    useRef,
    useState,
} from "react";
import { gameContext } from "../../contexts/gameContext";
import { nakamaClient } from "../../utils/nakama";
import { styled } from "@mui/material/styles";
import { useRouter } from "next/router";

const Root = styled("div")(({ theme }: { theme: Theme }) => ({
    width: "100%",
    overflow: "scroll",
    ...theme.typography.body2,
    "& > :not(style) + :not(style)": {
        marginTop: theme.spacing(2),
    },
}));

export const GameHistory = () => {
    const { shouldReconnect, infos, setShouldReconnect } =
        useContext(gameContext);
    const [matchNotFound, setMatchNotFound] = useState(false);
    const router = useRouter();
    const handleReconnect = useCallback(() => {
        nakamaClient.restoreSessionOrAuthenticate().then(() => {
            nakamaClient
                .reconnect()
                .then(() => {
                    setShouldReconnect(false);
                })
                .catch((err) => {
                    if (err.message === "Match not found") {
                        setMatchNotFound(true);
                        setTimeout(() => {
                            router.push("/");
                        }, 3000);
                    }
                });
        });
    }, []);

    const infoContainer = useRef(null);

    useEffect(() => {
        const scroll =
            infoContainer.current.scrollHeight -
            infoContainer.current.clientHeight;
        infoContainer.current.scrollTo(0, scroll);
    }, [infos]);

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
            {matchNotFound && "比赛已结束，即将返回首页"}
            <Root ref={infoContainer}>
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
