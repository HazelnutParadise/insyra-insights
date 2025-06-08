// ContextMenu 組件的類型定義

export interface ContextMenuItem {
    id: string;
    label?: string;
    icon?: string;
    danger?: boolean;
    type?: "separator" | "item";
    disabled?: boolean;
}

export interface ContextMenuConfig {
    [key: string]: ContextMenuItem[];
}

export interface ContextMenuActionEvent {
    action: string;
    context?: any;
}
