// 對話框服務 - 提供程式化調用對話框的功能
import { writable } from 'svelte/store';
import type { AlertOptions, ConfirmOptions, InputOptions } from '../types/dialog';

// 對話框狀態儲存
export const alertStore = writable<{
    visible: boolean;
    options: AlertOptions;
    resolve?: (value: void | PromiseLike<void>) => void;
}>({
    visible: false,
    options: { message: '' }
});

export const confirmStore = writable<{
    visible: boolean;
    options: ConfirmOptions;
    resolve?: (value: boolean | PromiseLike<boolean>) => void;
}>({
    visible: false,
    options: { message: '' }
});

export const inputStore = writable<{
    visible: boolean;
    options: InputOptions;
    resolve?: (value: string | null | PromiseLike<string | null>) => void;
}>({
    visible: false,
    options: { message: '' }
});

// Alert 函數 - 替代原生 alert
export function showAlert(messageOrOptions: string | AlertOptions): Promise<void> {
    return new Promise((resolve) => {
        const options: AlertOptions = typeof messageOrOptions === 'string'
            ? { message: messageOrOptions }
            : messageOrOptions;

        alertStore.set({
            visible: true,
            options,
            resolve
        });
    });
}

// Confirm 函數 - 替代原生 confirm
export function showConfirm(messageOrOptions: string | ConfirmOptions): Promise<boolean> {
    return new Promise((resolve) => {
        const options: ConfirmOptions = typeof messageOrOptions === 'string'
            ? { message: messageOrOptions }
            : messageOrOptions;

        confirmStore.set({
            visible: true,
            options,
            resolve
        });
    });
}

// Input 函數 - 替代原生 prompt
export function showInput(messageOrOptions: string | InputOptions): Promise<string | null> {
    return new Promise((resolve) => {
        const options: InputOptions = typeof messageOrOptions === 'string'
            ? { message: messageOrOptions }
            : messageOrOptions;

        inputStore.set({
            visible: true,
            options,
            resolve
        });
    });
}

// 關閉 Alert
export function closeAlert() {
    alertStore.update(state => {
        if (state.resolve) {
            state.resolve();
        }
        return {
            visible: false,
            options: { message: '' }
        };
    });
}

// 關閉 Confirm
export function closeConfirm(result: boolean = false) {
    confirmStore.update(state => {
        if (state.resolve) {
            state.resolve(result);
        }
        return {
            visible: false,
            options: { message: '' }
        };
    });
}

// 關閉 Input
export function closeInput(result: string | null = null) {
    inputStore.update(state => {
        if (state.resolve) {
            state.resolve(result);
        }
        return {
            visible: false,
            options: { message: '' }
        };
    });
}
