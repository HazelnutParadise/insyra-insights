// 對話框組件的類型定義

export interface DialogButton {
    id: string;
    label: string;
    type?: 'primary' | 'secondary' | 'danger';
    disabled?: boolean;
}

export interface AlertOptions {
    title?: string;
    message: string;
    buttonText?: string;
    type?: 'info' | 'warning' | 'error' | 'success';
}

export interface ConfirmOptions {
    title?: string;
    message: string;
    confirmText?: string;
    cancelText?: string;
    type?: 'info' | 'warning' | 'danger';
}

export interface InputOptions {
    title?: string;
    message: string;
    placeholder?: string;
    defaultValue?: string;
    confirmText?: string;
    cancelText?: string;
    type?: 'info' | 'warning' | 'danger';
    inputType?: 'text' | 'email' | 'password' | 'number';
}

export interface DialogEvent {
    action: string;
    result?: boolean;
}
