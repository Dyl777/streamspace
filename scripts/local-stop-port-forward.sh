#!/usr/bin/env bash
#
# local-stop-port-forward.sh - Stop all StreamSpace port forwards
#
# This script stops all running port forwards created by local-port-forward.sh
# by killing the background kubectl port-forward processes.
#

set -euo pipefail

# Colors for output
COLOR_RESET='\033[0m'
COLOR_BOLD='\033[1m'
COLOR_GREEN='\033[32m'
COLOR_YELLOW='\033[33m'
COLOR_BLUE='\033[34m'
COLOR_RED='\033[31m'

# Project configuration
PROJECT_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
PID_DIR="${PROJECT_ROOT}/.port-forward-pids"
LOG_DIR="${PROJECT_ROOT}/.port-forward-logs"

# Helper functions
log() {
    echo -e "${COLOR_BOLD}==>${COLOR_RESET} $*"
}

log_success() {
    echo -e "${COLOR_GREEN}✓${COLOR_RESET} $*"
}

log_error() {
    echo -e "${COLOR_RED}✗${COLOR_RESET} $*" >&2
}

log_info() {
    echo -e "${COLOR_BLUE}→${COLOR_RESET} $*"
}

log_warning() {
    echo -e "${COLOR_YELLOW}⚠${COLOR_RESET} $*"
}

# Stop all port forwards
stop_port_forwards() {
    log "Stopping port forwards..."
    echo ""

    if [ ! -d "${PID_DIR}" ]; then
        log_warning "No PID directory found (${PID_DIR})"
        log_info "No port forwards to stop"
        return 0
    fi

    local stopped=0
    local not_running=0

    for pid_file in "${PID_DIR}"/*.pid; do
        if [ -f "${pid_file}" ]; then
            local pid=$(cat "${pid_file}")
            local name=$(basename "${pid_file}" .pid)

            if kill -0 "${pid}" 2>/dev/null; then
                log_info "Stopping ${name} (PID: ${pid})..."
                kill "${pid}" 2>/dev/null || true

                # Wait for process to die
                local timeout=5
                local elapsed=0
                while kill -0 "${pid}" 2>/dev/null && [ $elapsed -lt $timeout ]; do
                    sleep 1
                    elapsed=$((elapsed + 1))
                done

                if kill -0 "${pid}" 2>/dev/null; then
                    log_warning "${name} did not stop gracefully, forcing..."
                    kill -9 "${pid}" 2>/dev/null || true
                fi

                log_success "${name} stopped"
                stopped=$((stopped + 1))
            else
                log_info "${name} was not running"
                not_running=$((not_running + 1))
            fi

            rm -f "${pid_file}"
        fi
    done

    echo ""
    if [ $stopped -gt 0 ]; then
        log_success "Stopped ${stopped} port forward(s)"
    fi

    if [ $not_running -gt 0 ]; then
        log_info "Cleaned up ${not_running} stale PID file(s)"
    fi

    if [ $stopped -eq 0 ] && [ $not_running -eq 0 ]; then
        log_warning "No port forwards were running"
    fi
}

# Clean up log directory
cleanup_logs() {
    if [ -d "${LOG_DIR}" ]; then
        log "Cleaning up log files..."
        rm -rf "${LOG_DIR}"
        log_success "Log files cleaned"
    fi
}

# Main execution
main() {
    echo -e "${COLOR_BOLD}═══════════════════════════════════════════════════${COLOR_RESET}"
    echo -e "${COLOR_BOLD}  Stop Port Forwards${COLOR_RESET}"
    echo -e "${COLOR_BOLD}═══════════════════════════════════════════════════${COLOR_RESET}"
    echo ""

    stop_port_forwards
    cleanup_logs

    echo ""
    echo -e "${COLOR_BOLD}═══════════════════════════════════════════════════${COLOR_RESET}"
    log_success "All port forwards stopped"
    echo -e "${COLOR_BOLD}═══════════════════════════════════════════════════${COLOR_RESET}"
    echo ""

    log_info "To restart port forwards:"
    echo "  ./scripts/local-port-forward.sh"
    echo ""
}

# Run main function
main "$@"
