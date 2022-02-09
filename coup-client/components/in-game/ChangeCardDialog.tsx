import React, { useContext, useState } from "react";
import { gameContext, ICard } from "../../contexts/gameContext";
import {
    Box,
    Button,
    Checkbox,
    Dialog,
    DialogTitle,
    ListItemText,
    MenuItem,
    OutlinedInput,
    Select,
    SelectChangeEvent,
} from "@mui/material";
import { nakamaClient } from "../../utils/nakama";
import { rolesMap } from "../../constants";

export interface ChangeCardDialogProps {
    open: boolean;
    handleClose: () => void;
}

export const ChangeCardDialog = (props: ChangeCardDialogProps) => {
    const { open, handleClose } = props;
    const { chooseCards, client } = useContext(gameContext);
    const [chosenCards, setChosenCards] = useState<string[]>([]);
    const allCards: ICard[] = client.cards.concat(chooseCards);
    const handleChange = (event: SelectChangeEvent<string[]>) => {
        const {
            target: { value },
        } = event;
        setChosenCards(value as string[]);
    };
    const handleSubmit = async () => {
        await nakamaClient.changeCard(chosenCards);
        handleClose();
    };
    return (
        <Dialog open={open} onBackdropClick={() => {}}>
            <DialogTitle sx={{ textAlign: "center" }}>
                选择要换的卡牌
            </DialogTitle>
            <Box
                sx={{
                    display: "flex",
                    flexDirection: "column",
                    alignItems: "center",
                    p: 1,
                    minWidth: 250,
                }}
            >
                <Select
                    sx={{
                        minWidth: 200,
                    }}
                    labelId="demo-multiple-checkbox-label"
                    id="demo-multiple-checkbox"
                    multiple
                    value={chosenCards}
                    onChange={handleChange}
                    input={<OutlinedInput label="已选" />}
                    renderValue={(selected) =>
                        selected
                            .map(
                                (cardId) =>
                                    rolesMap[
                                        allCards.find(
                                            (item) => item.id === cardId
                                        ).role
                                    ]
                            )
                            .join(" ")
                    }
                >
                    {allCards.map((card) => (
                        <MenuItem
                            disabled={
                                chosenCards.length === client.cards.length &&
                                !chosenCards.includes(card.id)
                            }
                            key={card.id}
                            value={card.id}
                        >
                            <Checkbox checked={chosenCards.includes(card.id)} />
                            <ListItemText primary={rolesMap[card.role]} />
                        </MenuItem>
                    ))}
                </Select>
                <Button
                    sx={{
                        mt: 2,
                    }}
                    onClick={() => handleSubmit()}
                    variant="contained"
                >
                    确认
                </Button>
            </Box>
        </Dialog>
    );
};
