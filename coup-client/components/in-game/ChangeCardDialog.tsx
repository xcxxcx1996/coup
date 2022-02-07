import React, { useContext } from "react";
import { gameContext, ICard } from "../../contexts/gameContext";
import {
    Box,
    Dialog,
    DialogTitle,
    MenuItem,
    OutlinedInput,
    Select,
    SelectChangeEvent,
} from "@mui/material";
import { nakamaClient } from "../../utils/nakama";
import { AbilityBtn, AbilityProps, IMenuItem } from "./ControlPanel";
import { rolesMap } from "../../constants";

export interface ChangeCardDialogProps {
    open: boolean;
    handleClose: () => void;
}

export const ChangeCardDialog = (props: ChangeCardDialogProps) => {
    const { open, handleClose } = props;
    const { cards, chooseCards, setCards } = useContext(gameContext);
    const allCards: ICard[] = cards.concat(chooseCards);
    const handleChange = (event: SelectChangeEvent<typeof cards>) => {
        const {
            target: { value },
        } = event;
        setCards(value as ICard[]);
    };
    return (
        <Dialog onClose={handleClose} open={open}>
            <DialogTitle sx={{ textAlign: "center" }}>选择换的卡牌</DialogTitle>
            <Box
                sx={{
                    display: "flex",
                    flexDirection: "column",
                    alignItems: "center",
                    p: 1,
                }}
            >
                <Select
                    labelId="demo-multiple-name-label"
                    id="demo-multiple-name"
                    multiple
                    value={cards}
                    onChange={handleChange}
                    input={<OutlinedInput label="Name" />}
                >
                    {allCards.map((card) => (
                        <MenuItem key={card.id} value={card.role}>
                            {rolesMap[card.role]}
                        </MenuItem>
                    ))}
                </Select>
            </Box>
        </Dialog>
    );
};
