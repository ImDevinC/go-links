import { ContentCopy, Delete, Visibility, WarningRounded } from "@mui/icons-material";
import { Alert, Button, DialogActions, DialogContent, DialogTitle, Divider, Modal, ModalDialog, Stack, Typography } from "@mui/joy";
import React, { useState } from "react";
import { LinkData, disableLink } from "../../services/api/links";
import aveta from 'aveta';

export interface LinkProps {
    link: LinkData
    onLinkChanged: () => void
}

export const Link = (props: LinkProps) => {
    const [openDialog, setOpenDialog] = useState(false)
    const [message, setMessage] = useState('')

    const copyToClipboard = () => {
        navigator.clipboard.writeText(`go/${props.link.name}`)
    }

    const handleDisableLink = async () => {
        const result = await disableLink(props.link.name)
        if (result.error) {
            setMessage(result.error)
            return
        }
        props.onLinkChanged()
    }

    return (
        <Stack direction="column" spacing={2} my={2}>
            <Stack direction="row" spacing={2} my={2}>
                <Stack direction="column" spacing={1} flexGrow={1}>
                    <Typography level="body-md">go/{props.link.name}</Typography>
                    <Typography level="body-sm">{props.link.description}</Typography>
                </Stack>
                <Stack direction="row" spacing={1} sx={{ alignItems: 'center' }}>
                    <Typography level="body-md" startDecorator={<Visibility />}>{aveta(props.link.views!, { precision: 2, lowercase: true })}</Typography>
                    <Button variant="plain" sx={{ mx: 0, p: 1 }} onClick={() => setOpenDialog(true)}><Delete /></Button>
                    <Button variant="plain" sx={{ mx: 0, p: 1 }} onClick={copyToClipboard}><ContentCopy /></Button>
                </Stack>
            </Stack>
            <Divider />
            <Modal open={openDialog} onClose={() => setOpenDialog(false)}>
                <ModalDialog variant="outlined" role="alertdialog">
                    <DialogTitle>
                        <WarningRounded />
                        Confirmation
                    </DialogTitle>
                    <Divider />
                    <DialogContent>
                        Are you sure you want to disable this link? This will impact all users.
                        {message !== '' ? <Alert color="danger">{message}</Alert> : null}
                    </DialogContent>
                    <DialogActions>
                        <Button variant="solid" color="danger" onClick={handleDisableLink}>Disable</Button>
                        <Button variant="plain" color="neutral" onClick={() => setOpenDialog(false)}>Cancel</Button>
                    </DialogActions>
                </ModalDialog>
            </Modal>
        </Stack >
    )
}