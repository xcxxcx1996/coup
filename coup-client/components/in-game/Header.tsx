import { Box } from "@mui/system";
import { useContext } from "react";
import { gameContext } from "../../contexts/gameContext";

export function Header() {
    const { currentPlayer, timeLeft } = useContext(gameContext);
    return (
        <Box
            sx={{
                p: 1,
                m: 1,
                display: "flex",
                justifyContent: "space-between",
                color: "red",
            }}
        >
            <div>{currentPlayer}的回合</div>
            {/*<div>回合时间: {timeLeft}s</div>*/}
        </Box>
    );
}
