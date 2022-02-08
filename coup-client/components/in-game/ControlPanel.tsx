import {
    Box,
    Button,
    Dialog,
    DialogTitle,
    Menu,
    MenuItem,
} from "@mui/material";
import React, { useContext, useState } from "react";
import KeyboardArrowDownIcon from "@mui/icons-material/KeyboardArrowDown";
import { gameContext } from "../../contexts/gameContext";
import { nakamaClient } from "../../utils/nakama";
import { rolesMap } from "../../constants";
import { State } from "../../constants/state";

export const isStateIdle = (playerState: number): boolean => {
    return playerState === State.IDLE;
};

export const isStateQuestion = (playerState: number): boolean => {
    return playerState === State.QUESTION;
};

export const isStateDenyMoney = (playerState: number): boolean => {
    return playerState === State.DENY_MONEY;
};

export const isStateDenyAssassin = (playerState: number): boolean => {
    return playerState === State.DENY_ASSASSIN;
};

export const isStateDenySteal = (playerState: number): boolean => {
    return playerState === State.DENY_STEAL;
};

export const isStateDiscard = (playerState: number): boolean => {
    return playerState === State.DISCARD;
};

export const isStateStart = (playerState: number): boolean => {
    return playerState === State.START;
};

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
                {items.map((item) => (
                    <MenuItem
                        dense={true}
                        divider={true}
                        key={item.text}
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
    const clientState = client.state;
    const isIdle = isStateIdle(clientState);
    const isDenyMoney = isStateDenyMoney(clientState);
    const isDenyAssassin = isStateDenyAssassin(clientState);
    const isDenySteal = isStateDenySteal(clientState);
    const handleClick = (ability: () => Promise<void>) => {
        return async () => {
            await ability();
            handleClose();
        };
    };
    return (
        <Dialog onClose={handleClose} open={open}>
            <DialogTitle sx={{ textAlign: "center" }}>
                选择要使用的卡牌技能
            </DialogTitle>
            <Box
                sx={{
                    display: "flex",
                    flexDirection: "column",
                    alignItems: "center",
                    p: 1,
                }}
            >
                <AbilityBtn
                    text="女王（防刺杀）"
                    onClick={handleClick(() => nakamaClient.denyKill(true))}
                    disabled={isIdle || !isDenyAssassin}
                />
                <AbilityBtn
                    text="男爵（收3金币）"
                    onClick={handleClick(nakamaClient.drawThreeCoins)}
                    disabled={isIdle}
                />
                <AbilityBtn
                    text="男爵（阻止收2金币）"
                    onClick={handleClick(() => nakamaClient.denyMoney(true))}
                    disabled={isIdle || !isDenyMoney}
                />
                <AbilityBtn
                    text="大使（换牌）"
                    onClick={handleClick(nakamaClient.changeCard)}
                    disabled={isIdle}
                />
                <AbilityBtn
                    text="大使（阻止偷金币）"
                    onClick={handleClick(() => nakamaClient.denySteal(true))}
                    disabled={isIdle || !isDenySteal}
                />
                <MenuButton
                    text={"队长（偷2金币）"}
                    items={users.map((u) => ({
                        text: u.name,
                        onClick: handleClick(() =>
                            nakamaClient.stealCoins(u.id)
                        ),
                    }))}
                    btnWidth={250}
                    menuItemWidth={250}
                    disabled={isIdle}
                />
                <AbilityBtn
                    text="队长（阻止偷金币）"
                    onClick={handleClick(() => nakamaClient.denySteal(true))}
                    disabled={isIdle || !isDenySteal}
                />
                <MenuButton
                    text={"刺客（刺杀）"}
                    items={users.map((u) => ({
                        text: u.name,
                        onClick: handleClick(() => nakamaClient.assassin(u.id)),
                    }))}
                    btnWidth={250}
                    menuItemWidth={250}
                    disabled={isIdle}
                />
            </Box>
        </Dialog>
    );
};

export const ControlPanel = () => {
    const [open, setOpen] = useState(false);
    const { users, client } = useContext(gameContext);
    const { cards, state } = client;
    const isIdle = isStateIdle(state);
    const isQuestion = isStateQuestion(state);
    const isDiscard = isStateDiscard(state);
    const isStart = isStateStart(state);
    const handleClose = () => {
        setOpen(false);
    };
    const handleClickOpen = () => {
        setOpen(true);
    };
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
                    { text: "1金币", onClick: () => nakamaClient.drawCoins(1) },
                    { text: "2金币", onClick: () => nakamaClient.drawCoins(2) },
                ]}
                text="获取金币"
                btnWidth={140}
                disabled={isIdle || !isStart}
            />
            <Button
                sx={{ width: "140px", m: 1 }}
                variant="contained"
                onClick={handleClickOpen}
                disabled={isIdle || !isStart}
            >
                使用技能
            </Button>
            <MenuButton
                items={[
                    {
                        text: "质疑",
                        onClick: () => nakamaClient.question(true),
                    },
                    {
                        text: "不质疑",
                        onClick: () => nakamaClient.question(false),
                    },
                ]}
                text="质疑/不质疑"
                btnWidth={140}
                disabled={isIdle || !isQuestion}
            />
            <MenuButton
                btnWidth={140}
                text={"政变"}
                items={users
                    .filter((u) => u.id !== nakamaClient?.session?.user_id)
                    .map((u) => ({
                        text: u.name,
                        onClick: () => nakamaClient.coup(u.id),
                    }))}
                disabled={isIdle || client.coins < 7 || !isStart}
            />
            <MenuButton
                btnWidth={140}
                text={"弃牌"}
                disabled={!isDiscard || isIdle}
                items={cards.map((c) => ({
                    text: rolesMap[c.role],
                    onClick: () => nakamaClient.discardCard(c.id),
                }))}
            />
            <AbilityDialog open={open} handleClose={handleClose} />
        </Box>
    );
};
