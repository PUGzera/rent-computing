
function FormError({message}) {
    return (
        <div style={{
            backgroundColor: '#ffe5e5',
            color: '#d8000c',
            border: '1px solid #f5c6cb',
            padding: '15px',
            borderRadius: '6px',
            marginTop: '15px',
            display: 'flex',
            alignItems: 'center',
            boxShadow: '0 2px 5px rgba(0,0,0,0.1)',
            display: message ? 'contents' : 'none'
            }}>
            <span style={{ marginRight: '10px' }}>⚠️</span>
            <span><strong>Error:</strong> {message}</span>
        </div>
    )
}

export default FormError