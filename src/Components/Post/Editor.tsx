import { useState } from 'react';
import ReactQuill from 'react-quill';
import 'react-quill/dist/quill.snow.css';

function EditorComponent() {
    const modules = {
        toolbar: [
            [{ header: '1' }, { header: '2' }, { font: [] }],
            [{ size: [] }],
            ['bold', 'italic', 'underline', 'strike', 'blockquote'],
            [
                { list: 'ordered' },
                { list: 'bullet' },
                { indent: '-1' },
                { indent: '+1' }
            ],
            ['link', 'image', 'video'],
            ['clean']
        ],
        clipboard: {
            // toggle to add extra line breaks when pasting HTML:
            matchVisual: false
        }
    };
    const [value, setValue] = useState('');

    return (
        <>
            <ReactQuill
                theme="snow"
                modules={modules}
                value={value}
                onChange={setValue}
                placeholder="What are your thoughts?"
            />
        </>
    );
}

export default EditorComponent;
