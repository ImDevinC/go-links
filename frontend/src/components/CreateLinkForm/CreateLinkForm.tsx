import { Alert, Box, Button, FormControl, FormLabel, IconButton, Input, Stack } from "@mui/joy";
import React, { FormEvent, useState } from "react";
import { LinkData, createLink } from "../../services/api/links";
import { CloseRounded } from "@mui/icons-material";

export const CreateLinkForm = () => {
    const [formData, setFormData] = useState<LinkData>(
        { url: '', name: '', description: '' }
    )
    const [linkSubmitted, setLinkSubmitted] = useState(false)
    const [responseMessage, setResponseMessage] = useState('')
    const [messageSeverity, setMessageSeverity] = useState<'success' | 'danger'>('success')
    const [showAlert, setShowAlert] = useState(false)

    const handleInputChange = (event: React.ChangeEvent<HTMLInputElement>) => {
        setFormData({ ...formData, [event.target.name]: event.target.value })
    }

    const handleSubmit = async (event: FormEvent<HTMLFormElement>) => {
        event.preventDefault();
        setLinkSubmitted(true);
        const resp = await createLink(formData);
        setLinkSubmitted(false);
        setShowAlert(true);
        if (resp.error) {
            setMessageSeverity('danger');
            setResponseMessage(resp.error);
            return;
        }
        setMessageSeverity('success');
        setResponseMessage(`Link go/${formData.name} created successfully!`);
        setFormData({ url: '', name: '', description: '' });
    }

    const disableAlert = () => setShowAlert(false);

    return (
        <Box sx={{ flex: 1, width: '100%' }}>
            <form onSubmit={handleSubmit}>
                <Stack spacing={4} sx={{
                    display: 'flex',
                    maxWidth: '800px',
                    mx: 'auto',
                    px: { xs: 2, md: 6 },
                    py: { xs: 2, md: 3 }
                }}>
                    <FormControl>
                        <FormLabel>Destination URL</FormLabel>
                        <Input value={formData.url} placeholder="https://google.com" name="url" required onChange={handleInputChange} />
                    </FormControl>
                    <FormControl>
                        <FormLabel>Golink Name (without go/ prefix)</FormLabel>
                        <Input value={formData.name} placeholder="google" required name="name" onChange={handleInputChange} />
                    </FormControl>
                    <FormControl>
                        <FormLabel>Description</FormLabel>
                        <Input value={formData.description} placeholder="Google Search Engine" name="description" required onChange={handleInputChange} />
                    </FormControl>
                    <Box sx={{ flexGrow: 1 }}>
                        <Button type="submit" disabled={linkSubmitted}>Create Link</Button>
                    </Box>
                    {
                        showAlert ?
                            <Alert color={messageSeverity} endDecorator={
                                <IconButton variant="plain" size="sm" color="neutral">
                                    <CloseRounded onClick={disableAlert} />
                                </IconButton>
                            }>{responseMessage}</Alert>
                            : null
                    }

                </Stack>
            </form>
        </Box>
    )
}