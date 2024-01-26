import { Alert, Box, Button, FormControl, FormLabel, IconButton, Input, Stack } from "@mui/joy";
import React, { FormEvent, useState } from "react";
import { LinkData, createLink } from "../../services/api/links";
import { CloseRounded } from "@mui/icons-material";

interface CreateLinkFormProps {
    onSuccess: () => void
}

export const CreateLinkForm = (props: CreateLinkFormProps) => {
    const [formData, setFormData] = useState<LinkData>(
        { url: '', name: '', description: '' }
    )
    const [loading, setLoading] = useState(false)
    const [errorMessage, setErrorMessage] = useState('')

    const handleInputChange = (event: React.ChangeEvent<HTMLInputElement>) => {
        setFormData({ ...formData, [event.target.name]: event.target.value })
    }

    const handleSubmit = async (event: FormEvent<HTMLFormElement>) => {
        event.preventDefault();
        setLoading(true);
        setErrorMessage('');
        const resp = await createLink(formData);
        setLoading(false);
        if (resp.error) {
            setErrorMessage(resp.error);
            return;
        }
        setErrorMessage('');
        setFormData({ url: '', name: '', description: '' });
        props.onSuccess();
    }

    const disableAlert = () => setErrorMessage('');

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
                        <Button type="submit" loading={loading}>Create Link</Button>
                    </Box>
                    {
                        errorMessage !== '' ?
                            <Alert color="danger" endDecorator={
                                <IconButton variant="plain" size="sm" color="neutral">
                                    <CloseRounded onClick={disableAlert} />
                                </IconButton>
                            }>{errorMessage}</Alert>
                            : null
                    }

                </Stack>
            </form>
        </Box>
    )
}