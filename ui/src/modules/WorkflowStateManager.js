import entityRelationships from 'modules/entityRelationships';

// Returns true if stack provided makes sense
export function isStackValid(stack) {
    if (stack.length < 2) return true;

    // stack is invalid when the stack is in one of three states:
    //
    // 1) entity -> (entity parent list) -> entity parent -> plus one entity -> nav away
    // 2) entity -> (entity matches list) -> match entity -> plus one -> nav away
    // 3) entity -> (entity contains-inferred list) -> contains-inferred entity -> plus one -> nav away

    let isParentState;
    let isMatchState;
    let isInferredState;

    stack.forEach(({ t: type }, i) => {
        if (i > 0 && i !== stack.length - 1) {
            const { t: prevType } = stack[i - 1];
            if (!isParentState) {
                isParentState = entityRelationships.isParent(type, prevType);
            }
            if (!isMatchState) {
                isMatchState = entityRelationships.isMatch(type, prevType);
            }
            if (!isInferredState) {
                isInferredState = entityRelationships.isContainedInferred(prevType, type);
            }
        }
        return false;
    });
    return !isParentState && !isMatchState && !isInferredState;
}

// Resets the current state based on minimal parameters
function baseStateStack(entityType, entityId) {
    const pageObj = {
        t: entityType
    };
    if (entityId) {
        pageObj.i = entityId;
    }

    return [pageObj];
}

// Checks state stack for overflow state/invalid state and returns a valid trimmed version
function trimStack(stack) {
    // Navigate away if:
    // If there's no more "room" in the stack

    // if the top entity is a parent of the entity before that then navigate away
    // List navigates to: Top single -> selected list
    // Entity navigates to : Entity page (maybe not)
    if (isStackValid(stack)) return stack;
    const { t: type, i: id } = stack.slice(-1)[0];
    if (!id) {
        const { t, i } = stack.slice(-2)[0];
        return [...baseStateStack(t, i), { t: type }];
    }
    return baseStateStack(type, id);
}

/**
 * Summary: Class that ensures the shape of a WorkflowState object
 * {
 *   useCase: 'text',
 *   stateStack: [{t: 'entityType', i: 'entityId'},{t: 'entityType', i: 'entityId'}]
 * }
 */
export class WorkflowState {
    constructor(useCase, stateStack) {
        this.useCase = useCase;
        this.stateStack = stateStack || [];
    }
}

export default class WorkflowStateMgr {
    constructor(workflowState, searchState) {
        if (workflowState) {
            const { useCase, stateStack } = workflowState;
            this.workflowState = new WorkflowState(useCase, stateStack);
        } else {
            this.workflowState = new WorkflowState();
        }
        this.searchState = { ...searchState };
    }

    // Resets the current state based on minimal parameters
    base(entityType, entityId, useCase) {
        const newUseCase = useCase || this.workflowState.useCase;
        const newStateStack = baseStateStack(entityType, entityId);

        this.workflowState = new WorkflowState(newUseCase, newStateStack);
        return this;
    }

    // Adds a list of entityType related to the current workflowState
    pushList(type) {
        const listState = { t: type };

        // if coming from dashboard
        if (!this.workflowState.stateStack.length) {
            this.workflowState.stateStack = [listState];
            return this;
        }

        const currentItem = this.workflowState.stateStack.slice(-1)[0];
        if (currentItem.t && !currentItem.i) {
            // replace the list type
            currentItem.t = type;
            return this;
        }

        this.workflowState.stateStack = trimStack([...this.workflowState.stateStack, listState]);
        return this;
    }

    // Selects an item in a list by Id
    pushListItem(id) {
        const currentItem = this.workflowState.stateStack.slice(-1)[0];
        currentItem.i = id;
        return this;
    }

    // Shows an entity in relation to the top entity in the workflow
    pushRelatedEntity(type, id) {
        const currentItem = this.workflowState.stateStack.slice(-1)[0];
        if (!currentItem.i)
            throw new Error(`Can't push related entity onto a list. Use pushListItem(id) instead.`);

        const newStack = trimStack([...this.workflowState.stateStack, { t: type, i: id }]);
        this.workflowState.stateStack = newStack;

        return this;
    }

    // Goes back one level to the nearest valid state
    pop() {
        if (this.workflowState.stateStack.length === 1)
            // A state stack has to have at least one item in it
            return this;

        this.workflowState.stateStack.pop();
        return this;
    }
}
