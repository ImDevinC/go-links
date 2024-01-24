import { Box, Button, Modal, ModalDialog, Typography } from "@mui/joy";
import React, { useState } from "react";
import { CreateLinkForm } from "../CreateLinkForm";

export const NewLinkButton = () => {
    const [open, setOpen] = useState(false)

    return (
        <Box alignItems="end" alignContent="end">
            <Button onClick={() => setOpen(true)}>Create</Button>
            <Modal open={open} onClose={() => setOpen(false)}>
                <ModalDialog aria-labelledby="nested-modal-title">
                    <Typography id="nested-modal-title" level="h2">Create new golink</Typography>
                    <CreateLinkForm onSuccess={() => setOpen(false)} />
                </ModalDialog>
            </Modal>
        </Box >
    )
}