import { Box } from "@mui/system";

export function Header() {
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
            <div>John&apos;s turn</div>
            <div>Time Left: 33s</div>
        </Box>
    );
}
