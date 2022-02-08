import { Box } from "@mui/system";
import { useContext } from "react";
import { gameContext } from "../../contexts/gameContext";

export function Header() {
    const { currentPlayer, timeLeft, client } = useContext(gameContext);
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
            <div>
                {client.id === currentPlayer.id ? "你" : currentPlayer.name}
                的回合
            </div>
            <div>回合时间: {timeLeft}s</div>
        </Box>
    );
}
