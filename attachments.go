package monzo

type Attachment struct {
	ID            string `json:"id"`
	Created       string `json:"created"`
	UserID        string `json:"user_id"`
	TransactionID string `json:"external_id"`
	FileURL       string `json:"file_url"`
	FileType      string `json:"file_type"`
}

// TODO untested
func (cl *Client) RegisterAttachment(transactionID, fileURL, fileType string) (*Attachment, error) {
	args := map[string]string{
		"external_id": transactionID,
		"file_url":    fileURL,
		"file_type":   fileType,
	}
	attachment := &Attachment{}
	if err := cl.request("POST", "/attachment/register", args, attachment); err != nil {
		return nil, err
	}
	return attachment, nil
}

type UploadURL struct {
	FileName string `json:"file_name"`
	FileType string `json:"file_type"`
}

// TODO untested
func (cl *Client) AttachmentUpload(fileName, fileType string) (*UploadURL, error) {
	args := map[string]string{
		"file_name": fileName,
		"file_type": fileType,
	}
	rsp := &UploadURL{}
	if err := cl.request("POST", "/attachment/register", args, rsp); err != nil {
		return nil, err
	}
	return rsp, nil
}
