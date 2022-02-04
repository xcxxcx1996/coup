import { Container } from "@mui/material";
import { Header } from "../../components/in-game/Header";
import { UserCarousel } from "../../components/in-game/UserCarousel";
import { GameHistory } from "../../components/in-game/GameHistory";
import { ControlPanel } from "../../components/in-game/ControlPanel";

function InGame() {
    return (
        <Container
            maxWidth="xs"
            sx={{
                background:
                    "-webkit-gradient(linear, left top, right top, color-stop(0%,rgba(233,223,196,1)), color-stop(1%,rgba(233,223,196,1)), color-stop(2%,rgba(237,227,200,1)), color-stop(24%,rgba(237,227,200,1)), color-stop(25%,rgba(235,221,195,1)), color-stop(48%,rgba(233,223,196,1)), color-stop(49%,rgba(235,221,195,1)), color-stop(52%,rgba(230,216,189,1)), color-stop(53%,rgba(230,216,189,1)), color-stop(54%,rgba(233,219,192,1)), color-stop(55%,rgba(230,216,189,1)), color-stop(56%,rgba(230,216,189,1)), color-stop(57%,rgba(233,219,192,1)), color-stop(58%,rgba(230,216,189,1)), color-stop(73%,rgba(230,216,189,1)), color-stop(74%,rgba(233,219,192,1)), color-stop(98%,rgba(233,219,192,1)), color-stop(100%,rgba(235,221,195,1)))",
                minHeight: "100vh",
            }}
        >
            <Header />
            <UserCarousel />
            <GameHistory />
            <ControlPanel />
        </Container>
    );
}

export default InGame;
