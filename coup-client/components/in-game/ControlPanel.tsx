import {
    Box,
    Button,
    Dialog,
    DialogTitle,
    Menu,
    MenuItem,
} from "@mui/material";
import React, { useContext, useEffect, useState } from "react";
import KeyboardArrowDownIcon from "@mui/icons-material/KeyboardArrowDown";
import { gameContext } from "../../contexts/gameContext";
import { nakamaClient } from "../../utils/nakama";
import { rolesMap } from "../../constants";
import { IUser } from "./UserCarousel";
import { ChangeCardDialog } from "./ChangeCardDialog";
import { useRouter } from "next/router";
import { useAbility } from "../../hooks/useAbility";
import { isStateChooseCard } from "../../utils/countState";
import { useControlPanel } from "../../hooks/useControlPanel";

export interface IMenuItem {
    text: string;
    onClick: () => void;
}

export interface MenuButtonProps {
    text: string;
    items: IMenuItem[];
    btnWidth: number;
    menuItemWidth?: number;
    disabled?: boolean;
}

const MenuButton = (props: MenuButtonProps) => {
    const {
        text,
        items,
        btnWidth,
        menuItemWidth = 140,
        disabled = false,
    } = props;
    const [anchorEl, setAnchorEl] = React.useState<null | HTMLElement>(null);
    const open = Boolean(anchorEl);
    const handleClick = (event: React.MouseEvent<HTMLButtonElement>) => {
        setAnchorEl(event.currentTarget);
    };
    const handleClose = () => {
        setAnchorEl(null);
    };
    return (
        <div>
            <Button
                sx={{ width: `${btnWidth}px`, m: 1 }}
                variant="contained"
                onClick={handleClick}
                endIcon={<KeyboardArrowDownIcon />}
                disableElevation={true}
                disabled={disabled}
            >
                {text}
            </Button>
            <Menu
                id="basic-menu"
                anchorEl={anchorEl}
                open={open}
                onClose={handleClose}
                MenuListProps={{
                    "aria-labelledby": "basic-button",
                }}
            >
                {items.map((item, index) => (
                    <MenuItem
                        dense={true}
                        divider={true}
                        key={index}
                        onClick={async () => {
                            await item.onClick();
                            handleClose();
                        }}
                        sx={{
                            width: `${menuItemWidth}px`,
                            display: "flex",
                            justifyContent: "center",
                        }}
                    >
                        {item.text}
                    </MenuItem>
                ))}
            </Menu>
        </div>
    );
};

export interface AbilityProps {
    open: boolean;
    handleClose: () => void;
}

export interface AbilityBtnProps {
    text: string;
    onClick: () => void;
    disabled?: boolean;
}

export const AbilityBtn = ({
    text,
    onClick,
    disabled = false,
}: AbilityBtnProps) => {
    return (
        <Button
            variant="contained"
            sx={{ width: "250px", m: 1 }}
            onClick={onClick}
            disabled={disabled}
        >
            {text}
        </Button>
    );
};

export const AbilityDialog = (props: AbilityProps) => {
    const { open, handleClose } = props;
    const { users, client } = useContext(gameContext);
    const targetUsers = users.filter((user: IUser) => user.id !== client.id);
    const disabled = useAbility();
    const handleClick = (ability: () => Promise<void>) => {
        return async () => {
            await ability();
            handleClose();
        };
    };
    return (
        <Dialog onClose={handleClose} open={open}>
            <DialogTitle sx={{ textAlign: "center" }}>
                ??????????????????????????????
            </DialogTitle>
            <Box
                sx={{
                    display: "flex",
                    flexDirection: "column",
                    alignItems: "center",
                    p: 1,
                }}
            >
                <MenuButton
                    text={"?????????????????????"}
                    items={[
                        {
                            text: "??????",
                            onClick: handleClick(() =>
                                nakamaClient.denyKill(true)
                            ),
                        },
                        {
                            text: "?????????",
                            onClick: handleClick(() =>
                                nakamaClient.denyKill(false)
                            ),
                        },
                    ]}
                    btnWidth={250}
                    menuItemWidth={250}
                    disabled={disabled.denyAssassin}
                />
                <AbilityBtn
                    text="????????????3?????????"
                    onClick={handleClick(nakamaClient.drawThreeCoins)}
                    disabled={disabled.drawThreeCoins}
                />
                <MenuButton
                    text={"??????????????????2?????????"}
                    items={[
                        {
                            text: "??????",
                            onClick: handleClick(() =>
                                nakamaClient.denyMoney(true)
                            ),
                        },
                        {
                            text: "?????????",
                            onClick: handleClick(() =>
                                nakamaClient.denyMoney(false)
                            ),
                        },
                    ]}
                    btnWidth={250}
                    menuItemWidth={250}
                    disabled={disabled.denyMoney}
                />
                <AbilityBtn
                    text="??????????????????"
                    onClick={handleClick(nakamaClient.changeCard)}
                    disabled={disabled.changeCard}
                />
                <MenuButton
                    text={"???????????????????????????"}
                    items={[
                        {
                            text: "??????",
                            onClick: handleClick(() =>
                                nakamaClient.denySteal(true, "DIPLOMAT")
                            ),
                        },
                        {
                            text: "?????????",
                            onClick: handleClick(() =>
                                nakamaClient.denySteal(false, "DIPLOMAT")
                            ),
                        },
                    ]}
                    btnWidth={250}
                    menuItemWidth={250}
                    disabled={disabled.denySteal}
                />
                <MenuButton
                    text={"????????????2?????????"}
                    items={targetUsers.map((u) => ({
                        text: u.name,
                        onClick: handleClick(() =>
                            nakamaClient.stealCoins(u.id)
                        ),
                    }))}
                    btnWidth={250}
                    menuItemWidth={250}
                    disabled={disabled.stealCoins}
                />
                <MenuButton
                    text={"???????????????????????????"}
                    items={[
                        {
                            text: "??????",
                            onClick: handleClick(() =>
                                nakamaClient.denySteal(true, "CAPTAIN")
                            ),
                        },
                        {
                            text: "?????????",
                            onClick: handleClick(() =>
                                nakamaClient.denySteal(false, "CAPTAIN")
                            ),
                        },
                    ]}
                    btnWidth={250}
                    menuItemWidth={250}
                    disabled={disabled.denySteal}
                />
                <MenuButton
                    text={"??????????????????"}
                    items={targetUsers.map((u) => ({
                        text: u.name,
                        onClick: handleClick(() => nakamaClient.assassin(u.id)),
                    }))}
                    btnWidth={250}
                    menuItemWidth={250}
                    disabled={disabled.assassin}
                />
            </Box>
        </Dialog>
    );
};

export const ControlPanel = () => {
    const [open, setOpen] = useState(false);
    const [openChooseCards, setOpenChooseCards] = useState(false);
    const { users, client, gameEnd } = useContext(gameContext);
    const { cards, state } = client;
    const isChooseCard = isStateChooseCard(state);
    const disabled = useControlPanel();

    const router = useRouter();
    const handleClose = () => {
        setOpen(false);
    };
    const handleClickOpen = () => {
        setOpen(true);
    };

    useEffect(() => {
        if (isChooseCard) {
            setOpenChooseCards(true);
        }
    }, [isChooseCard]);

    useEffect(() => {
        let timeout: ReturnType<typeof setTimeout>;
        console.log("-> gameEnd", gameEnd);
        if (gameEnd) {
            nakamaClient
                .leaveMatch()
                .then(
                    () => (timeout = setTimeout(() => router.push("/"), 5000))
                );
        }
        return () => clearTimeout(timeout);
    }, [gameEnd]);

    return (
        <Box
            sx={{
                display: "flex",
                flexDirection: "row",
                flexWrap: "wrap",
                justifyContent: "space-between",
                m: 1,
            }}
        >
            <MenuButton
                items={[
                    { text: "1??????", onClick: () => nakamaClient.drawCoins(1) },
                    { text: "2??????", onClick: () => nakamaClient.drawCoins(2) },
                ]}
                text="????????????"
                btnWidth={140}
                disabled={disabled.drawCoins}
            />
            <Button
                sx={{ width: "140px", m: 1 }}
                variant="contained"
                onClick={handleClickOpen}
                disabled={disabled.useAbility}
            >
                ????????????
            </Button>
            <MenuButton
                items={[
                    {
                        text: "??????",
                        onClick: () => nakamaClient.question(true),
                    },
                    {
                        text: "?????????",
                        onClick: () => nakamaClient.question(false),
                    },
                ]}
                text="??????/?????????"
                btnWidth={140}
                disabled={disabled.question}
            />
            <MenuButton
                btnWidth={140}
                text={"??????"}
                items={users
                    .filter((u) => u.id !== nakamaClient?.session?.user_id)
                    .map((u) => ({
                        text: u.name,
                        onClick: () => nakamaClient.coup(u.id),
                    }))}
                disabled={disabled.coup}
            />
            <MenuButton
                btnWidth={140}
                text={"??????"}
                disabled={disabled.discard}
                items={cards.map((c) => ({
                    text: rolesMap[c.role],
                    onClick: () => nakamaClient.discardCard(c.id),
                }))}
            />
            <ChangeCardDialog
                open={openChooseCards}
                handleClose={() => {
                    setOpenChooseCards(false);
                }}
            />
            <AbilityDialog open={open} handleClose={handleClose} />
        </Box>
    );
};
